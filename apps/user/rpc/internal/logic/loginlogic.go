package logic

import (
	"OuterIM/apps/user/models"
	"OuterIM/pkg/ctxdata"
	"OuterIM/pkg/encrypt"
	"OuterIM/pkg/xerr"
	"context"
	"github.com/pkg/errors"
	"time"

	"OuterIM/apps/user/rpc/internal/svc"
	"OuterIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "该手机号尚未注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {

	// 查找该号码是否已注册
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	//验证密码是否正确
	if !encrypt.ValidatePassword(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// token生成
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(
		l.svcCtx.Config.Jwt.AccessSecret,
		now,
		l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id,
	)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get json web token err %v", err)
	}
	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
