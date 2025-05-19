package handler

import (
	common "mscoin-common"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type WithdrawHandler struct {
	svcCtx *svc.ServiceContext
}

func NewWithdrawHandler(svcCtx *svc.ServiceContext) *WithdrawHandler {
	return &WithdrawHandler{svcCtx}
}

func (h *WithdrawHandler) QueryWithdrawCoin(w http.ResponseWriter, r *http.Request) {
	var req types.WithdrawReq
	l := logic.NewWithdrawLogic(r.Context(), h.svcCtx)
	resp, err := l.QueryWithdrawCoin(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func (h *WithdrawHandler) SendCode(w http.ResponseWriter, r *http.Request) {
	var req types.WithdrawReq
	l := logic.NewWithdrawLogic(r.Context(), h.svcCtx)
	resp, err := l.SendCode(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func (h *WithdrawHandler) WithdrawCode(w http.ResponseWriter, r *http.Request) {
	var req types.WithdrawReq
	if err := httpx.ParseForm(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	l := logic.NewWithdrawLogic(r.Context(), h.svcCtx)
	resp, err := l.WithdrawCode(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}

func (h *WithdrawHandler) Record(w http.ResponseWriter, r *http.Request) {
	var req types.WithdrawReq
	if err := httpx.ParseForm(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	l := logic.NewWithdrawLogic(r.Context(), h.svcCtx)
	resp, err := l.Record(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}
