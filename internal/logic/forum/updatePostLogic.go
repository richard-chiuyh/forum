package forumlogic

import (
	"context"
	"errors"
	"time"

	"forum/dao/db"
	"forum/dao/model"
	"forum/internal/logic"
	"forum/internal/svc"
	"forum/types/forum"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdatePost C 端会员仅能删除自己的内容（帖子或评论）
func (l *UpdatePostLogic) UpdatePost(in *forum.UpdatePostRequest) (*forum.UpdatePostResponse, error) {
	if in.Post == nil || in.Post.Id == 0 || in.Post.UserId == 0 {
		return &forum.UpdatePostResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}

	old, err := l.svcCtx.DBManager.PostDao.Get(l.ctx, in.Post.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &forum.UpdatePostResponse{BaseResp: logic.NoRecordErrorResponse()}, nil
		}
		l.Logger.Errorf("UpdatePost Get old content err(%+v)", err)
		return &forum.UpdatePostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	// 仅允许本人删除自己的内容
	if old.UserID != in.Post.UserId || in.Post.Status != forum.Status_STATUS_DELETED {
		return &forum.UpdatePostResponse{BaseResp: logic.NoPermissionResponse()}, nil
	}
	if old.Status == int32(forum.Status_STATUS_DELETED) {
		return &forum.UpdatePostResponse{BaseResp: logic.SuccessResponse()}, nil
	}

	old.Status = int32(forum.Status_STATUS_DELETED)
	if old.ParentID != 0 {
		// 评论删除：记录状态更新时间
		statusUpdateTime := time.Now().UnixMilli()
		old.StatusUpdateTime = &statusUpdateTime
	}
	if err := l.svcCtx.DBManager.PostDao.Update(l.ctx, old); err != nil {
		l.Logger.Errorf("UpdatePost update err(%+v)", err)
		return &forum.UpdatePostResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	if old.ParentID == 0 {
		// 帖子：从 ES 移除并广播帖子变更
		if err := l.svcCtx.ESManager.PostDao.Delete(l.ctx, old.ID); err != nil {
			l.Logger.Errorf("UpdatePost ES Delete err(%+v)", err)
		}
		// logic.SendWsPostChangeMessage(l.ctx, l.svcCtx.Kafka, old)
	} else {
		// 评论：父节点 comment_count - 1 并广播评论变更
		l.svcCtx.DBManager.Executor.Add(db.CountTask{
			Table:  model.Post{}.TableName(),
			Id:     old.ParentID,
			Column: db.COMMENT_COUNT,
			Delta:  -1,
		})
		// logic.SendWsforumCommentChangeMessage(l.ctx, l.svcCtx.Kafka, old)
	}

	return &forum.UpdatePostResponse{BaseResp: logic.SuccessResponse()}, nil
}
