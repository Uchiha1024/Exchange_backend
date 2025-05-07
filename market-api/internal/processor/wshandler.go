package processor

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-api/internal/model"
	"market-api/internal/ws"
)

type WebsocketHandler struct {
	wsServer *ws.WebsocketServer
}

func NewWebsocketHandler(wsServer *ws.WebsocketServer) *WebsocketHandler {
	return &WebsocketHandler{
		wsServer: wsServer,
	}
}

func (w *WebsocketHandler) HandleTradePlate(symbol string, plate *model.TradePlateResult) {
	bytes, _ := json.Marshal(plate)
	logx.Info("====买卖盘通知:", symbol, plate.Direction, ":", fmt.Sprintf("%d", len(plate.Items)))
	w.wsServer.BroadcastToNamespace("/", "/topic/market/trade-plate/"+symbol, string(bytes))
}
func (w *WebsocketHandler) HandleTrade(symbol string, data []byte) {
	//订单交易完成后 进入这里进行处理 订单就称为K线的一部分 数据量小 无法维持K线 K线来源 okx平台来
	//TODO implement me
	panic("implement me")
}

func (w *WebsocketHandler) HandleKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb) {
	logx.Info("================WebsocketHandler Start=======================")
	logx.Info("symbol:", symbol)
	thumb := thumbMap[symbol]
	if thumb == nil {
		thumb = kline.InitCoinThumb(symbol)
	}
	coinThumb := kline.ToCoinThumb(symbol, thumb)
	result := &model.CoinThumb{}
	copier.Copy(result, coinThumb)
	marshal, _ := json.Marshal(result)
	w.wsServer.BroadcastToNamespace("/", "/topic/market/thumb", string(marshal))

	bytes, _ := json.Marshal(kline)
	w.wsServer.BroadcastToNamespace("/", "/topic/market/kline/"+symbol, string(bytes))

	logx.Info("================WebsocketHandler End=======================")
}
