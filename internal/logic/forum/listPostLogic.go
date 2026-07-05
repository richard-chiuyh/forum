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
	"gorm.io/gorm"
)

type ListPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostLogic {
	return &ListPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListPost 列出内容。parent_id 语义：
//   - 0  : 帖子(根节点) feed
//   - >0 : 该节点的直接子节点(评论/回复)
//   - <0 : 仅评论(parent_id<>0)，通常配合 user_id 查询某用户的评论
func (l *ListPostLogic) ListPost(in *forum.ListPostRequest) (*forum.ListPostResponse, error) {
	in.Pagination = logic.PaginationCheck(in.Pagination)

	var (
		posts []*model.Post
		total int64
		err   error
	)

	// 内容搜索仅对帖子(根节点)生效，走 ES
	if in.ParentId == 0 && in.PostFilter != nil && in.PostFilter.Content != "" {
		posts, total, err = l.svcCtx.ESManager.PostDao.Search(l.ctx, in)
		if err != nil {
			l.Logger.Errorf("ListPost Search err(%+v)", err)
			return &forum.ListPostResponse{BaseResp: logic.ElasticsearchErrorResponse()}, err
		}
	} else {
		cond, scopes := l.buildScope(in)
		posts, total, err = l.svcCtx.DBManager.PostDao.List(l.ctx, cond, scopes...)
		if err != nil {
			l.Logger.Errorf("ListPost List err(%+v)", err)
			return &forum.ListPostResponse{BaseResp: logic.MysqlErrorResponse()}, err
		}
	}

	likedMap, favMap, err := l.getPostLikedAndFavMap(posts, in.OperatorId)
	if err != nil {
		l.Logger.Errorf("ListPost getPostLikedAndFavMap err(%+v)", err)
		return &forum.ListPostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	return &forum.ListPostResponse{
		BaseResp:   logic.SuccessResponse(),
		Posts:      dto.ToRpcPosts(posts, likedMap, favMap),
		Pagination: logic.PaginationResp(in.Pagination, total),
	}, nil
}

func (l *ListPostLogic) getPostLikedAndFavMap(posts []*model.Post, operatorId int64) (map[int64]bool, map[int64]bool, error) {
	if operatorId == 0 {
		return map[int64]bool{}, map[int64]bool{}, nil
	}
	ids := make([]int64, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, post.ID)
	}
	likedMap, err := logic.GetLikedMap(l.ctx, l.svcCtx.DBManager.LikeDao, ids, operatorId)
	if err != nil {
		return nil, nil, err
	}
	favMap, err := logic.GetFavMap(l.ctx, l.svcCtx.DBManager.FavDao, ids, operatorId)
	if err != nil {
		return nil, nil, err
	}
	return likedMap, favMap, nil
}

func (l *ListPostLogic) buildScope(in *forum.ListPostRequest) (cond *model.Post, scopes []db.ScopeF) {
	cond = &model.Post{
		UserID: in.UserId,
	}
	scopes = []db.ScopeF{
		statusScope(in.PostFilter),
		postFilterScope(in.PostFilter),
		db.SortScope(in.Sorts),
		db.OffsetAndLimitScope(logic.Offset(in.Pagination), logic.Limit(in.Pagination)),
	}
	switch {
	case in.ParentId > 0:
		// 列出某节点的直接子节点（评论/回复），并预加载父链以回填根帖子
		scopes = append(scopes, db.ParentIdEqScope(in.ParentId), preloadParentChainScope())
	case in.ParentId < 0:
		// 仅评论
		scopes = append(scopes, db.ParentIdNotZeroScope(), preloadParentChainScope())
	default:
		// 帖子(根节点)
		scopes = append(scopes, db.ParentIdEqScope(0))
	}
	return cond, scopes
}

// statusScope 默认仅返回已发布内容，除非显式指定状态筛选
func statusScope(filter *forum.PostFilter) db.ScopeF {
	return func(tx *gorm.DB) *gorm.DB {
		if filter != nil && len(filter.Status) > 0 {
			return tx.Where("status IN ?", filter.Status)
		}
		return tx.Where("status = ?", forum.Status_STATUS_PUBLISHED)
	}
}

func postFilterScope(filter *forum.PostFilter) db.ScopeF {
	return func(tx *gorm.DB) *gorm.DB {
		if filter == nil {
			return tx
		}
		if filter.PostId != 0 {
			tx = tx.Where("id = ?", filter.PostId)
		}
		if filter.Account != "" {
			tx = tx.Where("account like ?", "%"+filter.Account+"%")
		}
		if filter.PublishTimeStart != 0 {
			tx = tx.Where("publish_time >= ?", filter.PublishTimeStart)
		}
		if filter.PublishTimeEnd != 0 {
			tx = tx.Where("publish_time <= ?", filter.PublishTimeEnd)
		}
		return tx
	}
}

// preloadParentChainScope 预加载父链（≤2 跳）：Parent 为直接父，Parent.Parent 为根帖子（回复时）
func preloadParentChainScope() db.ScopeF {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Parent").Preload("Parent.Parent")
	}
}
