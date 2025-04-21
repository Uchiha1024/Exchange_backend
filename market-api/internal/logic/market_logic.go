package logic

import (
	"context"
	"grpc-common/market/types/market"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarketLogic) SymbolThumbTrend(req *types.MarketReq) ([]*types.CoinThumbResp, error) {
	// 声明一个变量用于存储币种缩略图数据列表
	var thrumbs []*market.CoinThumb
	// 从处理器中获取缓存的缩略图数据
	thrumb := l.svcCtx.Processor.GetThumb()
	// 标记是否成功从缓存获取数据
	isCache := false
	// 检查缓存数据是否存在
	if thrumb != nil {
		// 使用类型断言检查缓存数据的类型
		switch v := thrumb.(type) {
		case []*market.CoinThumb:
			// 如果类型匹配，将缓存数据转换为正确的类型
			thrumbs = v
			// 标记成功从缓存获取数据
			isCache = true
		}
	}

	// 如果缓存中没有有效数据，则通过RPC调用获取
	if !isCache {
		// 创建一个带超时的上下文，5秒后自动取消
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// 确保在函数返回时取消上下文
		defer cancel()

		// 调用RPC服务获取最新的币种缩略图趋势数据
		symbolThumbResp, err := l.svcCtx.MarketRpc.FindSymbolThumbTrend(ctx, &market.MarketReq{
			Ip: req.Ip, // 传递请求中的IP地址
		})
		// 检查RPC调用是否出错
		if err != nil {
			return nil, err
		}
		// 将RPC返回的数据保存到thrumbs变量中
		thrumbs = symbolThumbResp.List
	}

	// 声明一个变量用于存储转换后的响应数据
	var coinThumbResp []*types.CoinThumbResp
	// 使用copier库将market.CoinThumb类型转换为types.CoinThumbResp类型
	if err := copier.Copy(&coinThumbResp, thrumbs); err != nil {
		return nil, err
	}

	// 返回转换后的数据和nil错误
	return coinThumbResp, nil
}

func (l *MarketLogic) SymbolThumb(req *types.MarketReq) ([]*types.CoinThumbResp, error) {

	// 声明一个变量用于存储币种缩略图数据列表
	var thrumbs []*market.CoinThumb
	// 从处理器中获取缓存的缩略图数据
	thrumb := l.svcCtx.Processor.GetThumb()

	// 检查缓存数据是否存在
	if thrumb != nil {
		// 使用类型断言检查缓存数据的类型
		switch v := thrumb.(type) {
		case []*market.CoinThumb:
			// 如果类型匹配，将缓存数据转换为正确的类型
			thrumbs = v
		}
	}

	// 声明一个变量用于存储转换后的响应数据
	var coinThumbResp []*types.CoinThumbResp
	// 使用copier库将market.CoinThumb类型转换为types.CoinThumbResp类型
	if err := copier.Copy(&coinThumbResp, thrumbs); err != nil {
		return nil, err
	}

	// 返回转换后的数据和nil错误
	return coinThumbResp, nil
}

func (l *MarketLogic) SymbolInfo(req *types.MarketReq) (*types.ExchangeCoinResp, error) {

	// 创建一个带超时的上下文，5秒后自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 确保在函数返回时取消上下文
	defer cancel()

	symbolInfoResp, err := l.svcCtx.MarketRpc.FindSymbolInfo(ctx, &market.MarketReq{
		Ip:     req.Ip,
		Symbol: req.Symbol,
	})

	if err != nil {
		return nil, err
	}

	resp := types.ExchangeCoinResp{}

	if err := copier.Copy(&resp, symbolInfoResp); err != nil {
		return nil, err
	}

	return &resp, nil

}

func (l *MarketLogic) CoinInfo(req *types.MarketReq) (*types.Coin, error) {

	// 创建一个带超时的上下文，5秒后自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 确保在函数返回时取消上下文
	defer cancel()

	coinInfoResp, err := l.svcCtx.MarketRpc.FindCoinInfo(ctx, &market.MarketReq{
		Unit: req.Unit,
	})

	if err != nil {
		logx.Errorw("RPC-CoinInfo", logx.Field("err", err))
		return nil, err
	}

	coinInfo := types.Coin{}

	if err := copier.Copy(&coinInfo, coinInfoResp); err != nil {
		return nil, err
	}

	return &coinInfo, nil

}

func (l *MarketLogic) History(req *types.MarketReq) (*types.HistoryKline, error) {
	// 创建一个带超时的上下文，5秒后自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 确保在函数返回时取消上下文
	defer cancel()
	historyKline, err := l.svcCtx.MarketRpc.HistoryKline(ctx, &market.MarketReq{
		Symbol:     req.Symbol,
		From:       req.From,
		To:         req.To,
		Resolution: req.Resolution,
	})
	if err != nil {
		logx.Errorw("RPC-HistoryKline", logx.Field("err", err))
		return nil, err
	}

	histories := historyKline.List
	var list = make([][]any, len(histories))
	for i, v := range histories {
		content := make([]any, 6)
		content[0] = v.Time
		content[1] = v.Open
		content[2] = v.High
		content[3] = v.Low
		content[4] = v.Close
		content[5] = v.Volume
		list[i] = content

	}
	return &types.HistoryKline{
		List: list,
	}, nil
}
