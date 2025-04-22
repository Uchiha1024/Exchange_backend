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

type LoginHandler struct {
	svcCtx *svc.ServiceContext
}

func NewLoginHandler(svcCtx *svc.ServiceContext) *LoginHandler {
	return &LoginHandler{
		svcCtx: svcCtx,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginReq

	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	newResult := common.NewResult()

	req.Ip = tools.GetRemoteClientIp(r)

	l := logic.NewLoginLogic(r.Context(), h.svcCtx)
	resp, err := l.Login(&req)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}

func (h *LoginHandler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	newResult := common.NewResult()
	token := r.Header.Get("x-auth-token")
	
	l := logic.NewLoginLogic(r.Context(), h.svcCtx)
	resp, err := l.CheckLogin(token)
	result := newResult.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}

