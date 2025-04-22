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

type AssetHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAssetHandler(svcCtx *svc.ServiceContext) *AssetHandler {
	return &AssetHandler{
		svcCtx: svcCtx,
	}
}

func (h *AssetHandler) FindWalletBySymbol(w http.ResponseWriter, r *http.Request) {
	var req *types.AssetReq
	// 解析请求参数
	if err := httpx.ParsePath(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	ip := tools.GetRemoteClientIp(r)
	req.Ip = ip
	// 调用服务层
	logic := logic.NewAssetLogic(r.Context(), h.svcCtx)
	resp, err := logic.FindWalletBySymbol(req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}
