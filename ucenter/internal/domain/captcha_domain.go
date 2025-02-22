package domain

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"webCoin-common/tools"
)

type CaptchaDomain struct {
}

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

func (d *CaptchaDomain) Verity(server string, vid string, key string, token string, scene int, ip string) bool {
	//根据vaptcha的api调用文档，发送相关post请求既可
	resp, err := tools.Post(server, &vaptchaReq{
		Id:        vid,
		Secretkey: key,
		Scene:     scene,
		Token:     token,
		Ip:        ip,
	})
	if err != nil {
		logx.Errorf("CaptchaDomain Verify post err : %s", err.Error())
		return false
	}
	var vaptchaRsp *vaptchaRsp
	err = json.Unmarshal(resp, &vaptchaRsp)
	if err != nil {
		logx.Errorf("CaptchaDomain Verify Unmarshal respBytes err : %s", err.Error())
		return false
	}
	if vaptchaRsp != nil && vaptchaRsp.Success == 1 {
		logx.Info("CaptchaDomain Verify success")
		return true
	}
	return false
}

func NewCaptchaDomain() *CaptchaDomain {
	return &CaptchaDomain{}
}
