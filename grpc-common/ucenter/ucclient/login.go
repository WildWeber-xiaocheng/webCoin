// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.6
// Source: register.proto

package ucclient

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"grpc-common/ucenter/types/login"
)

type (
	LoginReq = login.LoginReq
	LoginRes = login.LoginRes

	Login interface {
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error)
	}

	defaultLogin struct {
		cli zrpc.Client
	}
)

func NewLogin(cli zrpc.Client) Login {
	return &defaultLogin{
		cli: cli,
	}
}

func (m *defaultLogin) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error) {
	client := login.NewLoginClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}
