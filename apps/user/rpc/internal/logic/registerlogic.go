package logic

import (
	"OuterIM/apps/user/models"
	"OuterIM/apps/user/rpc/internal/svc"
	"OuterIM/apps/user/rpc/user"
	"OuterIM/pkg/ctxdata"
	"OuterIM/pkg/encrypt"
	"OuterIM/pkg/wuid"
	"OuterIM/pkg/xerr"
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegistered = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号已被注册")
)

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

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {

	// 查找该号码，防止二次注册
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)

	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get user by phone while registry err %v , req %v", err, in.Phone)
	}

	if userEntity != nil {
		return nil, errors.WithStack(ErrPhoneIsRegistered)
	}

	//定义用户数据
	userEntity = &models.Users{
		Id:       wuid.GetUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 密码加密
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "password hash err %v", err)
		}
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	// 向数据库中Insert
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert new user err %v ", err)
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
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
