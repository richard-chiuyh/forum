package db

import (
	"context"

	"forum/dao/model"
	"forum/dbcore"

	"gorm.io/gorm"
)

type favOpr struct {
	dbCtx dbcore.Context
}

func newFavOpr(conn gorm.DB) FavDao {
	return &favOpr{
		dbCtx: dbcore.NewContext(&conn),
	}
}

func (o *favOpr) Create(ctx context.Context, fav *model.Fav) error {
	return o.dbCtx.DB(ctx).Create(fav).Error
}

func (o *favOpr) Get(ctx context.Context, userId, targetId int64) (fav *model.Fav, err error) {
	return fav, o.dbCtx.DB(ctx).
		Where("user_id = ? AND target_id = ?", userId, targetId).
		Unscoped().First(&fav).Error
}

func (o *favOpr) List(ctx context.Context, cond *model.Fav, scopes ...ScopeF) (favs []*model.Fav, total int64, err error) {
	return favs, total, o.dbCtx.DB(ctx).Where(cond).Scopes(scopes...).Find(&favs).Offset(-1).Limit(-1).Count(&total).Error
}

func (o *favOpr) Update(ctx context.Context, fav *model.Fav) error {
	return o.dbCtx.DB(ctx).Save(fav).Error
}
