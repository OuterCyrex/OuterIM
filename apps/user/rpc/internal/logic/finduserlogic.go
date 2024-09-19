package logic

import (
	"OuterIM/apps/user/models"
	"context"
	"github.com/jinzhu/copier"

	"OuterIM/apps/user/rpc/internal/svc"
	"OuterIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {

	var (
		userEntities []*models.Users
		err          error
	)

	if in.Phone != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntities = append(userEntities, userEntity)
		}
	} else if in.Name != "" {
		userEntities, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntities, err = l.svcCtx.UsersModel.ListById(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	var resp []*user.UserEntity
	_ = copier.Copy(&resp, userEntities)

	return &user.FindUserResp{
		User: resp,
	}, nil
}
