package logic

import (
	"OuterIM/pkg/xerr"
	"context"
	"github.com/jinzhu/copier"
	PkgErr "github.com/pkg/errors"

	"OuterIM/apps/social/rpc/internal/svc"
	"OuterIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	friendList, err := l.svcCtx.ListByUserid(l.ctx, in.UserId)
	if err != nil {
		return nil, PkgErr.Wrapf(xerr.NewDBErr(), "find userlist err %v req %v", err, in.UserId)
	}

	var resp []*social.Friends
	_ = copier.Copy(&resp, friendList)

	return &social.FriendListResp{
		List: resp,
	}, nil
}
