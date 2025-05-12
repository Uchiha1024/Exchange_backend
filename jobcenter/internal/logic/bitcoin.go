package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/ucclient"
	"jobcenter/internal/database"
	"jobcenter/internal/domain"
	"log"
	"mscoin-common/tools"
	"sync"
)

type BitCoinConfig struct {
	Address string
}
type BitCoin struct {
	wg            sync.WaitGroup
	ch            cache.Cache
	assetRpc      ucclient.Asset
	bitCoinDomain *domain.BitCoinDomain
	queueDomain   *domain.QueueDomain
}

// Do 扫描BTC交易 查找符合系统address的交易 进行存储
func (b *BitCoin) Do(address string) {
	b.wg.Add(1)
	go b.ScanTx(address)
	b.wg.Wait()

}

func (b *BitCoin) ScanTx(btcAddress string) {
	//1. redis查询是否有记录区块，获取到已处理的区块高度 dealBlocks
	var dealBlocksStr string
	b.ch.Get("BTC::TX", &dealBlocksStr)
	var dealBlocks int64
	if dealBlocksStr == "" {
		dealBlocks = 2428713
	} else {
		dealBlocks = tools.ToInt64(dealBlocksStr)
	}
	//2. 根据getmininginfo获取到现在的区块高度 currentBlocks
	currentBlocks, err := b.getMiningInfo(btcAddress)
	if err != nil {
		log.Println(err)
		b.wg.Done()
		return
	}
	//3. 根据currentBlocks-dealBlocks 如果小于等于0 不需要扫描
	diff := currentBlocks - dealBlocks
	if diff <= 0 {
		b.wg.Done()
		return
	}
	//4. 获取系统中的BTC的address列表
	ctx := context.Background()
	address, err := b.assetRpc.GetAddress(ctx, &asset.AssetReq{
		CoinName: "BTC",
	})
	if err != nil {
		log.Println(err)
		b.wg.Done()
		return
	}
	addressList := address.List
	//5. 循环 根据getblockhash 获取 blockhash
	for i := currentBlocks; i > dealBlocks; i-- {
		blockHash, err := b.getBlockHash(i, btcAddress)
		if err != nil {
			log.Println(err)
			b.wg.Done()
			continue
		}
		//6. 通过getblock 获取 交易id列表
		txIdList, err := b.getBlock(blockHash, btcAddress)
		if err != nil {
			log.Println(err)
			b.wg.Done()
			continue
		}
		//7. 循环交易id列表 获取到交易详情 得到 vout内容
		for _, txId := range txIdList {
			txResult, err := b.getRawTransaction(txId, btcAddress)
			if err != nil {
				log.Println(err)
				b.wg.Done()
				continue
			}

			inputAddressList := make([]string, len(txResult.Vin))
			for i, vin := range txResult.Vin {
				if vin.Txid == "" {
					continue
				}
				inputTx, err := b.getRawTransaction(vin.Txid, btcAddress)
				if err != nil {
					log.Println(err)
					b.wg.Done()
					continue
				}
				vout := inputTx.Vout[vin.Vout]
				inputAddressList[i] = vout.ScriptPubKey.Address
			}
			//8. 根据vout中的address和上方address列表进行匹配，如果匹配，我们认为是充值
			for _, vout := range txResult.Vout {
				voutAddress := vout.ScriptPubKey.Address
				flag := false
				//9. 做一个处理，根据vint的交易 查询input的address，
				//   如果address和vout当中和系统匹配的address一样，我们认为不是充值  2 0.5 1.5
				for _, inputAddress := range inputAddressList {
					if inputAddress != "" && voutAddress != "" && inputAddress == voutAddress {
						flag = true
					}
				}
				if flag {
					continue
				}
				for _, address := range addressList {
					if address != "" && address == voutAddress {
						//匹配上了 //10. 找到充值数据，存入mongo，同时发送kafka进行下一步处理（存入member_transaction表）
						//充值
						err := b.bitCoinDomain.Recharge(txResult.TxId, vout.Value, voutAddress, txResult.Time, txResult.Blockhash)
						if err != nil {
							log.Println(err)
							b.wg.Done()
							continue
						}
						//kafka处理
						b.queueDomain.SendRecharge(vout.Value, voutAddress, txResult.Time)
					}
				}

			}
		}

	}

	//11. 记录redis的区块高度
	b.ch.Set("BTC::TX", currentBlocks)
	b.wg.Done()
}

//	{
//	   "result": {
//	       "blocks": 2428737,
//	       "difficulty": 104649090.3850967,
//	       "networkhashps": 440403219589005.6,
//	       "pooledtx": 216,
//	       "chain": "test",
//	       "warnings": "Unknown new rules activated (versionbit 28)"
//	   },
//	   "error": null,
//	   "id": "mscoin"
//	}
type MiningInfoResult struct {
	Id     string     `json:"id"`
	Error  string     `json:"error"`
	Result MiningInfo `json:"result"`
}
type MiningInfo struct {
	Blocks        int     `json:"blocks"`
	Difficulty    float64 `json:"difficulty"`
	Networkhashps float64 `json:"networkhashps"`
	Pooledtx      int     `json:"pooledtx"`
	Chain         string  `json:"chain"`
	Warnings      string  `json:"warnings"`
}

func (b *BitCoin) getMiningInfo(addr string) (int64, error) {
	//{
	//    "jsonrpc": "1.0",
	//    "method": "getmininginfo",
	//    "params":[],
	//    "id": "mscoin"
	//}
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getmininginfo"
	params["params"] = []int{}
	params["id"] = "mscoin"
	headers := make(map[string]string)
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="
	bytes, err := tools.PostWithHeader(addr, params, headers, "")
	if err != nil {
		return 0, err
	}
	var result MiningInfoResult
	json.Unmarshal(bytes, &result)
	if result.Error != "" {
		return 0, errors.New(result.Error)
	}
	return int64(result.Result.Blocks), nil
}

type BlockHashResult struct {
	Id     string `json:"id"`
	Error  string `json:"error"`
	Result string `json:"result"`
}

func (b *BitCoin) getBlockHash(height int64, addr string) (string, error) {
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getblockhash"
	params["params"] = []int64{height}
	params["id"] = "mscoin"
	headers := make(map[string]string)
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="
	bytes, err := tools.PostWithHeader(addr, params, headers, "")
	if err != nil {
		return "", err
	}
	var result BlockHashResult
	json.Unmarshal(bytes, &result)
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.Result, nil
}

type BlockResult struct {
	Id     string      `json:"id"`
	Error  string      `json:"error"`
	Result BlockSimple `json:"result"`
}
type BlockSimple struct {
	Hash string   `json:"hash"`
	Tx   []string `json:"tx"`
	Time int64    `json:"time"`
}

func (b *BitCoin) getBlock(blockHash, addr string) ([]string, error) {
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getblock"
	params["params"] = []any{blockHash, 1}
	params["id"] = "mscoin"
	headers := make(map[string]string)
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="
	bytes, err := tools.PostWithHeader(addr, params, headers, "")
	if err != nil {
		return nil, err
	}
	var result BlockResult
	json.Unmarshal(bytes, &result)
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return result.Result.Tx, nil
}

type RawTransactionResult struct {
	Id     string         `json:"id"`
	Error  string         `json:"error"`
	Result RawTransaction `json:"result"`
}
type RawTransaction struct {
	TxId      string `json:"txid"`
	Hash      string `json:"hash"`
	Locktime  int64  `json:"locktime"`
	Version   int    `json:"version"`
	Size      int    `json:"size"`
	Vsize     int    `json:"vsize"`
	Weight    int    `json:"weight"`
	Vin       []Vin  `json:"vin"`
	Vout      []Vout `json:"vout"`
	Time      int64  `json:"time"`
	Hex       string `json:"hex"`
	Blocktime int64  `json:"blocktime"`
	Blockhash string `json:"blockhash"`
}

type Vin struct {
	Txid        string            `json:"txid"`
	Vout        int               `json:"vout"`
	Txinwitness []string          `json:"txinwitness"`
	Sequence    int64             `json:"sequence"`
	ScriptSig   map[string]string `json:"scriptSig"`
}

type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}
type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

func (b *BitCoin) getRawTransaction(txId, addr string) (*RawTransaction, error) {
	params := make(map[string]any)
	params["jsonrpc"] = "1.0"
	params["method"] = "getrawtransaction"
	params["params"] = []any{txId, true}
	params["id"] = "mscoin"
	headers := make(map[string]string)
	headers["Authorization"] = "Basic Yml0Y29pbjoxMjM0NTY="
	bytes, err := tools.PostWithHeader(addr, params, headers, "")
	if err != nil {
		return nil, err
	}
	var result RawTransactionResult
	json.Unmarshal(bytes, &result)
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return &result.Result, nil
}

func NewBitCoin(ch cache.Cache, assetRpc ucclient.Asset, client *database.MongoClient, kafkaClient *database.KafkaClient) *BitCoin {
	return &BitCoin{
		ch:            ch,
		assetRpc:      assetRpc,
		bitCoinDomain: domain.NewBitCoinDomain(client),
		queueDomain:   domain.NewQueueDomain(kafkaClient),
	}
}
