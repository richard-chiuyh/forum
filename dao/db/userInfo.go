package db

import (
	"context"

	"forum/dao/model"
	"forum/dbcore"

	"gorm.io/gorm"
)

type userInfoOpr struct {
	dbCtx dbcore.Context
}

func newUserInfoOpr(conn gorm.DB) UserInfoDao {
	return &userInfoOpr{
		dbCtx: dbcore.NewContext(&conn),
	}
}

func (c *userInfoOpr) Create(ctx context.Context, userInfo *model.UserInfo) error {
	return c.dbCtx.DB(ctx).Create(userInfo).Error
}

func (c *userInfoOpr) Update(ctx context.Context, userInfo *model.UserInfo) error {
	return c.dbCtx.DB(ctx).Save(userInfo).Error
}

func (c *userInfoOpr) GetByUserID(ctx context.Context, userId int64) (userInfo *model.UserInfo, err error) {
	return userInfo, c.dbCtx.DB(ctx).Where("user_id = ?", userId).First(&userInfo).Error
}
