package domain

import (
	"encoding/json"
	"mscoin-common/tools"

	"github.com/zeromicro/go-zero/core/logx"
)

type vaptchaReq struct {
	Id        string `json:"id"`
	Secretkey string `json:"secretkey"`
	Scene     int    `json:"scene"`
	Token     string `json:"token"`
	Ip        string `json:"ip"`
}
type vaptchaRsp struct {
	Success int    `json:"success"`
	Score   int    `json:"score"`
	Msg     string `json:"msg"`
}

type CaptchaDomain struct {
}

func NewCaptchaDomain() *CaptchaDomain {
	return &CaptchaDomain{}
}

func (d *CaptchaDomain) VerifyCaptcha(
	server string,
	vid string,
	key string,
	token string,
	scene int,
	ip string,
) bool {
	resp, err := tools.Post(server, &vaptchaReq{
		Id:        vid,
		Secretkey: key,
		Scene:     scene,
		Token:     token,
		Ip:        ip,
	})

	if err != nil {
		// 打印日志
		logx.Errorf("VerifyCaptcha error: %v", err)
		return false
	}

	result := &vaptchaRsp{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		logx.Errorf("VerifyCaptcha Unmarshal error: %v", err)
		return false
	}

	if result.Success != 1 {
		logx.Errorf("VerifyCaptcha failed, result: %v", result)
		return false
	}

	return result.Success == 1

}
