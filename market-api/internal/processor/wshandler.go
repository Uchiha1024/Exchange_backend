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



/* 
WebsocketHandler 实现了 MarketHandler 接口的所有方法：
HandleTrade
HandleKLine
HandleTradePlate
这就是为什么它可以被传入 AddHandler 方法。
在 Go 语言中，只要一个类型实现了某个接口的所有方法，它就自动实现了该接口，不需要显式声明。
这是 Go 语言的接口实现机制的一个特点，称为隐式接口实现。
这个设计的目的是：
通过 WebSocket 向客户端实时推送市场数据
处理不同类型的市场数据（K线、交易、盘口）
将数据广播到不同的 WebSocket 主题
这是一个典型的观察者模式实现，WebsocketHandler 作为一个观察者，监听和处理市场数据的变化，
并通过 WebSocket 推送给客户端。 
*/



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
	logx.Info("============coinThumbResult:", string(marshal))
	w.wsServer.BroadcastToNamespace("/", "/topic/market/thumb", string(marshal))

	bytes, _ := json.Marshal(kline)
	w.wsServer.BroadcastToNamespace("/", "/topic/market/kline/"+symbol, string(bytes))

	logx.Info("================WebsocketHandler End=======================")
}


