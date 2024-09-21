package logic

import (
	"context"

	"OuterIM/apps/social/rpc/internal/svc"
	"OuterIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInListLogic) GroupPutInList(in *social.GroupPutinListReq) (*social.GroupPutinResp, error) {
	// todo: add your logic here and delete this line

	return &social.GroupPutinResp{}, nil
}
