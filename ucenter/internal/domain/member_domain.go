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
	memberRepo repo.MemberRepo //repo小写，是防止外部直接操作repo，而不用domain
}

func (d *MemberDomain) FindByPhone(ctx context.Context, phone string) (*model.Member, error) {
	mem, err := d.memberRepo.FindByPhone(ctx, phone)
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
	mem.Avatar = "https://mszlu.oss-cn-beijing.aliyuncs.com/mscoin/defaultavatar.png"
	mem.Salt = salt
	err = d.memberRepo.Save(ctx, mem)
	if err != nil {
		logx.Error(err)
		return errors.New("数据库异常")
	}
	return nil
}

// 登录次数不重要，这里就没处理error
func (d *MemberDomain) UpdateLoginCount(ctx context.Context, id int64, step int) {
	err := d.memberRepo.UpdateLoginCount(ctx, id, step)
	if err != nil {
		logx.Error(err)
	}
}

func NewMemberDomain(db *msdb.MsDB) *MemberDomain {
	return &MemberDomain{
		memberRepo: dao.NewMemberDao(db),
	}
}
