package forumlogic

import (
	"context"

	"forum/dao/db"
	"forum/dao/model"
	"forum/dto"
	"forum/internal/logic"
	"forum/internal/svc"
	"forum/types/forum"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFavLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFavLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavLogic {
	return &ListFavLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取收藏列表
func (l *ListFavLogic) ListFav(in *forum.ListFavRequest) (*forum.ListFavResponse, error) {
	if in.UserId == 0 {
		return &forum.ListFavResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}
	in.Pagination = logic.PaginationCheck(in.Pagination)
	cond := &model.Fav{
		UserID: in.UserId,
		Status: int32(forum.Status_STATUS_PUBLISHED),
	}
	scopes := []db.ScopeF{
		db.TargetIdsInScope(in.TargetIds),
		favPreloadTargetChainScope(),
		db.OffsetAndLimitScope(logic.Offset(in.Pagination), logic.Limit(in.Pagination)),
		db.SortScope(in.Sorts),
	}
	favs, total, err := l.svcCtx.DBManager.FavDao.List(l.ctx, cond, scopes...)
	if err != nil {
		l.Logger.Errorf("ListFav List err(%+v)", err)
		return &forum.ListFavResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}
	likedMap, err := l.getFavLikedMap(favs, in.UserId)
	if err != nil {
		l.Logger.Errorf("ListFav getFavLikedMap err(%+v)", err)
		return &forum.ListFavResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}
	return &forum.ListFavResponse{
		BaseResp:   logic.SuccessResponse(),
		Pagination: logic.PaginationResp(in.Pagination, total),
		Favs:       dto.ToRpcFavs(favs, likedMap),
	}, nil
}

func (l *ListFavLogic) getFavLikedMap(favs []*model.Fav, operatorId int64) (map[int64]bool, error) {
	if operatorId == 0 {
		return map[int64]bool{}, nil
	}
	targetIds := make([]int64, 0, len(favs))
	for _, fav := range favs {
		targetIds = append(targetIds, fav.TargetID)
	}
	return logic.GetLikedMap(l.ctx, l.svcCtx.DBManager.LikeDao, targetIds, operatorId)
}

// favPreloadTargetChainScope 预加载收藏目标及其父链（评论目标需要父链以回填根帖子）
func favPreloadTargetChainScope() db.ScopeF {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Target").Preload("Target.Parent").Preload("Target.Parent.Parent")
	}
}
