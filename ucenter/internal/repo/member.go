package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberRepo interface {
	FindByPhone(ctx context.Context, phone string) (*model.Member, error)
}
