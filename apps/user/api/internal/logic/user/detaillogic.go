package user

import (
	"OuterIM/apps/user/rpc/user"
	"OuterIM/pkg/ctxdata"
	"context"
	"github.com/jinzhu/copier"

	"OuterIM/apps/user/api/internal/svc"
	"OuterIM/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {

	uid := ctxdata.GetUId(l.ctx)

	UserInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})

	if err != nil {
		return nil, err
	}

	var res types.User
	_ = copier.Copy(&res, UserInfoResp.User)
	return &types.UserInfoResp{Info: res}, nil
}