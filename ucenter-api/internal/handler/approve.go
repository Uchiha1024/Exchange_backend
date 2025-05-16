package handler

import (
	common "mscoin-common"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type ApproveHandler struct {
	svcCtx *svc.ServiceContext
}

func NewApproveHandler(svcCtx *svc.ServiceContext) *ApproveHandler {
	return &ApproveHandler{svcCtx}
}

func (h *ApproveHandler) SecuritySetting(w http.ResponseWriter, r *http.Request) {
	var req types.ApproveReq
	l := logic.NewApproveLogic(r.Context(), h.svcCtx)
	resp, err := l.FindSecuritySetting(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)

}
