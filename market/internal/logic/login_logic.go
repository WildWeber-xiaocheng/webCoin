package logic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/login"
	"time"
	"ucenter/internal/domain"
	"ucenter/internal/svc"
	"webCoin-common/tools"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	CaptchaDomain *domain.CaptchaDomain
	MemberDomain  *domain.MemberDomain
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		CaptchaDomain: domain.NewCaptchaDomain(),
		MemberDomain:  domain.NewMemberDomain(svcCtx.Db),
	}
}

func (l *LoginLogic) Login(in *login.LoginReq) (*login.LoginRes, error) {
	//logx.Info("ucenter rpc login by phone call...")
	//1.校验 人机验证是否通过
	//引入domain层，将各个具体的处理逻辑都抽象出来
	isVerify := l.CaptchaDomain.Verity(
		in.Captcha.Server,
		l.svcCtx.Config.Captcha.Vid,
		l.svcCtx.Config.Captcha.Key,
		in.Captcha.Token,
		2, //2 代表注册场景
		in.Ip)
	if !isVerify {
		logx.Error("人机校验未通过")
		return nil, errors.New("人机校验未通过")
	}
	logx.Info("人机校验通过....")

	//2.校验密码
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	member, err := l.MemberDomain.FindByPhone(ctx, in.GetUsername())
	if err != nil {
		logx.Error(err)
		return nil, errors.New("登录失败")
	}
	if member == nil {
		return nil, errors.New("此用户未注册")
	}
	password := member.Password
	salt := member.Salt
	verify := tools.Verify(in.Password, salt, password, nil)
	if !verify {
		return nil, errors.New("密码不正确")
	}

	//3.登录成功，jwt生成token，提供给前端，前端调用传递token，我们进行token认证既可
	secret := l.svcCtx.Config.JWT.AccessSecret
	expire := l.svcCtx.Config.JWT.AccessExpire
	token, err := l.getJwtToken(secret, time.Now().Unix(), expire, member.Id)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("未知错误，请联系管理员")
	}
	//4.返回登录信息
	loginCount := member.LoginCount + 1
	//登录次数不是太重要，放到协程处理
	go func() {
		l.MemberDomain.UpdateLoginCount(context.Background(), member.Id, 1)
	}()
	return &login.LoginRes{
		Token:         token,
		Id:            member.Id,
		Username:      member.Username,
		MemberLevel:   member.MemberLevelStr(),
		MemberRate:    member.MemberRate(),
		RealName:      member.RealName,
		Country:       member.Country,
		Avatar:        member.Avatar,
		PromotionCode: member.PromotionCode,
		SuperPartner:  member.SuperPartner,
		LoginCount:    int32(loginCount),
	}, nil
}

// 生成token
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
