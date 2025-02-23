package domain

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
	"webCoin-common/msdb"
	"webCoin-common/tools"
)

type MemberDomain struct {
	MemberRepo repo.MemberRepo
}

func (d *MemberDomain) FindByPhone(ctx context.Context, phone string) (*model.Member, error) {
	mem, err := d.MemberRepo.FindByPhone(ctx, phone)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("数据库异常")
	}
	return mem, nil
}

func (d *MemberDomain) Register(ctx context.Context, phone string, password string, username string, country string, partner string, promotion string) error {
	mem := model.NewMember()
	//password密码 需要进行md5加密,并且需要加盐，因为Md5加密不安全（可以通过彩虹表破解）
	//member表的字段比较多，并且所有的字段都不为Null，也就是很多字段要填写默认值，这里使用一个工具类，通过反射来填充默认值
	err := tools.Default(mem)
	if err != nil {
		logx.Error(err)
		return errors.New("赋默认值异常")
	}
	salt, pwd := tools.Encode(password, nil)
	mem.Username = username
	mem.Country = country
	mem.Password = pwd
	mem.MobilePhone = phone
	mem.FillSuperPartner(partner)
	mem.PromotionCode = promotion
	mem.MemberLevel = model.GENERAL
	mem.Salt = salt
	err = d.MemberRepo.Save(ctx, mem)
	if err != nil {
		logx.Error(err)
		return errors.New("数据库异常")
	}
	return nil
}

func NewMemberDomain(db *msdb.MsDB) *MemberDomain {
	return &MemberDomain{
		MemberRepo: dao.NewMemberDao(db),
	}
}
