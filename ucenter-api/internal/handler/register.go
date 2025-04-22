package handler

import (
	
	common "mscoin-common"
	"mscoin-common/tools"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type RegisterHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRegisterHandler(svcCtx *svc.ServiceContext) *RegisterHandler {
	return &RegisterHandler{
		svcCtx: svcCtx,
	}
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.Request

	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	newResult := common.NewResult()


	req.Ip = tools.GetRemoteClientIp(r)

	l := logic.NewRegisterLogic(r.Context(), h.svcCtx)
	resp, err := l.Register(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}

func (h *RegisterHandler) SendCode(w http.ResponseWriter, r *http.Request) {
	var req types.SendCodeRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	// 先创建对象
	l := logic.NewRegisterLogic(r.Context(), h.svcCtx)
	resp, err := l.SendCode(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}
