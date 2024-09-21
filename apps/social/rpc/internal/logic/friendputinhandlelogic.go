package logic

import (
	"OuterIM/apps/social/socialmodels"
	"OuterIM/pkg/constants"
	"OuterIM/pkg/xerr"
	"context"
	PkgErr "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"OuterIM/apps/social/rpc/internal/svc"
	"OuterIM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请已通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsg("好友申请已拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 获取好友申请记录

	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, PkgErr.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err %v req %v", err, in.FriendReqId)
	}

	// 验证是否有处理
	switch constants.HandlerResult(friendReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, PkgErr.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, PkgErr.WithStack(ErrFriendReqBeforeRefuse)
	}

	friendReq.HandleResult.Int64 = int64(in.HandleResult)
	// 修改申请结果
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendReq); err != nil {
			return PkgErr.Wrapf(xerr.NewDBErr(), "update friend requset err %v,req %v", err, friendReq)
		}

		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		friends := []*socialmodels.Friends{
			{
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			},
			{
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}

		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return PkgErr.Wrapf(xerr.NewDBErr(), "friends inserts err %v , req %v", err, friendReq)
		}
		return nil
	},
	)

	return &social.FriendPutInHandleResp{}, nil
}
