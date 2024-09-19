package user

import (
	"OuterIM/apps/user/rpc/user"
	"context"
	"github.com/jinzhu/copier"

	"OuterIM/apps/user/api/internal/svc"
	"OuterIM/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {

	RegisterResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})

	if err != nil {
		return nil, err
	}

	var res types.RegisterResp
	_ = copier.Copy(&res, RegisterResp)
	return &res, nil
}
