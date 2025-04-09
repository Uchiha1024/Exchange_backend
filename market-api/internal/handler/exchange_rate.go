package handler

import (
	"market-api/internal/logic"
	"market-api/internal/svc"
	"market-api/internal/types"
	common "mscoin-common"
	"mscoin-common/tools"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type ExchangeRateHandler struct {
	svcCtx *svc.ServiceContext
}

func NewExchangeRateHandler(svcCtx *svc.ServiceContext) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		svcCtx: svcCtx,
	}
}

func (h *ExchangeRateHandler) UsdRate(w http.ResponseWriter, r *http.Request) {
	var req types.RateRequest
	if err := httpx.ParsePath(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	newResult := common.NewResult()
	//获取一下ip
	req.Ip = tools.GetRemoteClientIp(r)
	l := logic.NewExchangeRateLogic(r.Context(), h.svcCtx)
	resp, err := l.UsdRate(&req)
	result := newResult.Deal(resp.Rate, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}
