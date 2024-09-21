package logic

import (
	"OuterIM/apps/social/socialmodels"
	"OuterIM/pkg/constants"
	"OuterIM/pkg/xerr"
	"context"
	"database/sql"
	"errors"
	"time"

	"OuterIM/apps/social/rpc/internal/svc"
	"OuterIM/apps/social/rpc/social"

	PkgErr "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {

	// friend already
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, PkgErr.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, err
	}

	// already pull request
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, PkgErr.Wrapf(xerr.NewDBErr(), "find friendsRequest by uid and fid err %v req %v", err, in)
	}
	if friendReqs != nil {
		return &social.FriendPutInResp{}, err
	}
	// new pull
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	if err != nil {
		return nil, PkgErr.Wrapf(xerr.NewDBErr(), "find friendRequest by rid and uid err %v req %v", err, in)
	}
	return &social.FriendPutInResp{}, nil
}
