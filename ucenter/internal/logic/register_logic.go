package logic

import (
	"context"
	"errors"
	"grpc-common/ucenter/types/register"
	"time"
	"ucenter/internal/svc"
	"webCoin-common/tools"

	"github.com/zeromicro/go-zero/core/logx"
)

const RegisterCacheKey = "REGISTER::"

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) RegisterByPhone(in *register.RegReq) (*register.RegRes, error) {
	logx.Info("ucenter rpc register by phone call...")
	return &register.RegRes{}, nil
}

func (l *RegisterLogic) SendCode(in *register.CodeReq) (*register.NoRes, error) {
	//1.生成验证码
	code := tools.Rand4Num()
	//2.发送验证码，可使用第三方api来实现
	go func() {
		logx.Info("调用短信平台发送验证码成功")
	}()
	logx.Infof("验证码为：%v", code)
	//3.验证码存入redis
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := l.svcCtx.Cache.SetWithExpireCtx(ctx, RegisterCacheKey+in.Phone, code, 15*time.Minute)
	if err != nil {
		return nil, errors.New("验证法发送失败")
	}
	//不能return nil,nil 因为如果返回是nil，grpc不能进行序列化
	return &register.NoRes{}, nil
}
