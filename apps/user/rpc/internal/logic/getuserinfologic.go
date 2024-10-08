package logic

import (
	"OuterIM/apps/user/models"
	"context"
	"errors"
	"github.com/jinzhu/copier"

	"OuterIM/apps/user/rpc/internal/svc"
	"OuterIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var resp user.UserEntity
	_ = copier.Copy(&resp, userEntity)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
