package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/login"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"webCoin-common/tools"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	loginReq := &login.LoginReq{}
	if err := copier.Copy(loginReq, req); err != nil {
		return nil, err
	}
	loginRes, err := l.svcCtx.UCLoginRpc.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}
	resp = &types.LoginRes{}
	if err := copier.Copy(resp, loginRes); err != nil {
		return nil, err
	}
	return
}

func (l *LoginLogic) CheckLogin(token string) (bool, error) {
	//没有用rpc，因为比较简单，直接api层进行处理
	_, err := tools.ParseToken(token, l.svcCtx.Config.JWT.AccessSecret)
	if err != nil {
		logx.Error(err)
		return false, nil
	}
	return true, nil
}
