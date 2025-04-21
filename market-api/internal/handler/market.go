package handler

import (
	"market-api/internal/logic"
	"market-api/internal/svc"
	"market-api/internal/types"
	common "mscoin-common"
	"mscoin-common/tools"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type MarketHandler struct {
	svcCtx *svc.ServiceContext
}

func NewMarketHandler(svcCtx *svc.ServiceContext) *MarketHandler {
	return &MarketHandler{
		svcCtx: svcCtx,
	}
}

func (h *MarketHandler) SymbolThumbTrend(w http.ResponseWriter, r *http.Request) {
	var req types.MarketReq
	newResult := common.NewResult()

	req.Ip = tools.GetRemoteClientIp(r)
	l := logic.NewMarketLogic(r.Context(), h.svcCtx)
	resp, err := l.SymbolThumbTrend(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}

func (h *MarketHandler) SymbolThumb(w http.ResponseWriter, r *http.Request) {
	var req types.MarketReq
	newResult := common.NewResult()

	req.Ip = tools.GetRemoteClientIp(r)
	l := logic.NewMarketLogic(r.Context(), h.svcCtx)
	resp, err := l.SymbolThumb(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func (h *MarketHandler) SymbolInfo(w http.ResponseWriter, r *http.Request) {
	var req types.MarketReq
	if err := httpx.ParseForm(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	newResult := common.NewResult()

	req.Ip = tools.GetRemoteClientIp(r)
	l := logic.NewMarketLogic(r.Context(), h.svcCtx)
	resp, err := l.SymbolInfo(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}

func (h *MarketHandler) CoinInfo(w http.ResponseWriter, r *http.Request) {

	var req types.MarketReq
	if err := httpx.ParseForm(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	newResult := common.NewResult()

	req.Ip = tools.GetRemoteClientIp(r)
	l := logic.NewMarketLogic(r.Context(), h.svcCtx)
	logx.Infof("CoinInfo req: %v", req)
	resp, err := l.CoinInfo(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}


func (h *MarketHandler) History(w http.ResponseWriter, r *http.Request) {
	var req types.MarketReq
	if err := httpx.ParseForm(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	newResult := common.NewResult()

	req.Ip = tools.GetRemoteClientIp(r)
	l := logic.NewMarketLogic(r.Context(), h.svcCtx)
	logx.Infof("CoinInfo req: %v", req)
	resp, err := l.History(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}
