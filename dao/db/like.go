package db

import (
	"context"

	"forum/dao/model"
	"forum/dbcore"

	"gorm.io/gorm"
)

type likeOpr struct {
	dbCtx dbcore.Context
}

func newLikeOpr(conn gorm.DB) LikeDao {
	return &likeOpr{
		dbCtx: dbcore.NewContext(&conn),
	}
}

func (o *likeOpr) Create(ctx context.Context, like *model.Like) error {
	return o.dbCtx.DB(ctx).Create(like).Error
}

// Get 通过 (user_id, target_id) 查找；包含已软删除/状态为已删除的记录，方便 toggle 复用历史行
func (o *likeOpr) Get(ctx context.Context, userId, targetId int64) (like *model.Like, err error) {
	return like, o.dbCtx.DB(ctx).
		Where("user_id = ? AND target_id = ?", userId, targetId).
		Unscoped().First(&like).Error
}

func (o *likeOpr) List(ctx context.Context, cond *model.Like, scopes ...ScopeF) (likes []*model.Like, total int64, err error) {
	return likes, total, o.dbCtx.DB(ctx).Where(cond).Scopes(scopes...).Find(&likes).Offset(-1).Limit(-1).Count(&total).Error
}

func (o *likeOpr) Update(ctx context.Context, like *model.Like) error {
	return o.dbCtx.DB(ctx).Save(like).Error
}
