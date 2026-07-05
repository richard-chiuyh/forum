package forumlogic

import (
	"context"
	"errors"

	"forum/dao/db"
	"forum/dao/model"
	"forum/internal/logic"
	"forum/internal/svc"
	"forum/types/forum"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Like 切换点赞状态：未点赞→点赞，已点赞→取消，已取消→重新点赞
func (l *LikeLogic) Like(in *forum.LikeRequest) (*forum.LikeResponse, error) {
	if in.UserId == 0 || in.TargetId == 0 {
		return &forum.LikeResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}

	// 校验目标内容存在
	if _, err := l.svcCtx.DBManager.PostDao.Get(l.ctx, in.TargetId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &forum.LikeResponse{BaseResp: logic.NoRecordErrorResponse()}, nil
		}
		l.Logger.Errorf("Like get target err(%+v)", err)
		return &forum.LikeResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	like, err := l.svcCtx.DBManager.LikeDao.Get(l.ctx, in.UserId, in.TargetId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("Like LikeDao.Get err(%+v)", err)
		return &forum.LikeResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	delta := int64(1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		like = &model.Like{
			TargetID:  in.TargetId,
			UserID:    in.UserId,
			Status:    int32(forum.Status_STATUS_PUBLISHED),
			IPAddress: in.IpAddress,
		}
		err = l.svcCtx.DBManager.LikeDao.Create(l.ctx, like)
	} else {
		// 严格 toggle：只在已发布/已取消两种状态间切换
		if like.Status == int32(forum.Status_STATUS_PUBLISHED) {
			like.Status = int32(forum.Status_STATUS_DELETED)
			delta = -1
		} else {
			like.Status = int32(forum.Status_STATUS_PUBLISHED)
		}
		if in.IpAddress != "" {
			like.IPAddress = in.IpAddress
		}
		err = l.svcCtx.DBManager.LikeDao.Update(l.ctx, like)
	}
	if err != nil {
		l.Logger.Errorf("Like create/update err(%+v)", err)
		return &forum.LikeResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	// 计数：目标内容 like_count ± delta
	l.svcCtx.DBManager.Executor.Add(db.CountTask{
		Table:  model.Post{}.TableName(),
		Id:     in.TargetId,
		Column: db.LIKE_COUNT,
		Delta:  delta,
	})

	return &forum.LikeResponse{BaseResp: logic.SuccessResponse()}, nil
}
