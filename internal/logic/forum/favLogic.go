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

type FavLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavLogic {
	return &FavLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Fav 切换收藏状态
func (l *FavLogic) Fav(in *forum.FavRequest) (*forum.FavResponse, error) {
	if in.UserId == 0 || in.TargetId == 0 {
		return &forum.FavResponse{BaseResp: logic.ParamErrorResponse()}, nil
	}

	// 校验目标内容存在
	if _, err := l.svcCtx.DBManager.PostDao.Get(l.ctx, in.TargetId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &forum.FavResponse{BaseResp: logic.NoRecordErrorResponse()}, nil
		}
		l.Logger.Errorf("Fav get target err(%+v)", err)
		return &forum.FavResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	fav, err := l.svcCtx.DBManager.FavDao.Get(l.ctx, in.UserId, in.TargetId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("Fav FavDao.Get err(%+v)", err)
		return &forum.FavResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	delta := int64(1)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fav = &model.Fav{
			TargetID: in.TargetId,
			UserID:   in.UserId,
			Status:   int32(forum.Status_STATUS_PUBLISHED),
		}
		err = l.svcCtx.DBManager.FavDao.Create(l.ctx, fav)
	} else {
		if fav.Status == int32(forum.Status_STATUS_PUBLISHED) {
			fav.Status = int32(forum.Status_STATUS_DELETED)
			delta = -1
		} else {
			fav.Status = int32(forum.Status_STATUS_PUBLISHED)
		}
		err = l.svcCtx.DBManager.FavDao.Update(l.ctx, fav)
	}
	if err != nil {
		l.Logger.Errorf("Fav create/update err(%+v)", err)
		return &forum.FavResponse{BaseResp: logic.MysqlErrorResponse()}, err
	}

	// 计数：目标内容 fav_count ± delta
	l.svcCtx.DBManager.Executor.Add(db.CountTask{
		Table:  model.Post{}.TableName(),
		Id:     in.TargetId,
		Column: db.FAV_COUNT,
		Delta:  delta,
	})

	return &forum.FavResponse{BaseResp: logic.SuccessResponse()}, nil
}
