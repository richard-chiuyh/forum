package forumlogic

import (
	"context"

	"forum/dao/db"
	"forum/dao/model"
	"forum/dto"
	"forum/internal/logic"
	"forum/internal/svc"
	"forum/types/forum"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLikeLogic {
	return &ListLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListLike 列出当前用户的点赞记录（仅返回 status=正常 的）
func (l *ListLikeLogic) ListLike(in *forum.ListLikeRequest) (*forum.ListLikeResponse, error) {
	if in.UserId == 0 {
		return &forum.ListLikeResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}
	in.Pagination = logic.PaginationCheck(in.Pagination)
	cond := &model.Like{
		UserID: in.UserId,
		Status: int32(forum.Status_STATUS_PUBLISHED),
	}
	scopes := []db.ScopeF{
		db.TargetIdsInScope(in.TargetIds),
		db.OffsetAndLimitScope(logic.Offset(in.Pagination), logic.Limit(in.Pagination)),
		db.SortScope(in.Sorts),
	}
	likes, total, err := l.svcCtx.DBManager.LikeDao.List(l.ctx, cond, scopes...)
	if err != nil {
		l.Logger.Errorf("ListLike List err(%+v)", err)
		return &forum.ListLikeResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}
	return &forum.ListLikeResponse{
		BaseResp:   logic.SuccessResponse(),
		Pagination: logic.PaginationResp(in.Pagination, total),
		Likes:      dto.ToRpcLikes(likes),
	}, nil
}
