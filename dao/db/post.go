package db

import (
	"context"

	"forum/dao/model"
	"forum/dbcore"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type postOpr struct {
	dbCtx dbcore.Context
}

func newPostOpr(conn gorm.DB) PostDao {
	return &postOpr{
		dbCtx: dbcore.NewContext(&conn),
	}
}

func (p *postOpr) Create(ctx context.Context, post *model.Post) error {
	return p.dbCtx.DB(ctx).Create(post).Error
}

func (p *postOpr) Update(ctx context.Context, post *model.Post) error {
	logx.WithContext(ctx).Infof("Update Post: %+v", post)
	return p.dbCtx.DB(ctx).Save(post).Error
}

func (p *postOpr) Get(ctx context.Context, id int64) (post *model.Post, err error) {
	return post, p.dbCtx.DB(ctx).Where("id = ?", id).First(&post).Error
}

func (p *postOpr) List(ctx context.Context, cond *model.Post, scopes ...ScopeF) (posts []*model.Post, total int64, err error) {
	return posts, total, p.dbCtx.DB(ctx).Where(cond).Scopes(scopes...).Find(&posts).Offset(-1).Limit(-1).Count(&total).Error
}
