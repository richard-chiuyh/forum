package model

import "gorm.io/plugin/soft_delete"

const TableNameCommUserInfo = "comm_user_info"

// UserInfo 社区用户信息表
type UserInfo struct {
	ID           int64                 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       int64                 `gorm:"column:user_id;not null;default:0;uniqueIndex:uniq_user_id" json:"user_id"`
	Account      string                `gorm:"column:account;type:varchar(32);not null;default:''" json:"account"`
	LastSeenPost int64                 `gorm:"column:last_seen_post;not null;default:0" json:"last_seen_post"`
	PostCount    int32                 `gorm:"column:post_count;not null;default:0" json:"post_count"`
	CreateTime   int64                 `gorm:"column:create_time;not null;autoCreateTime:milli" json:"create_time"`
	UpdateTime   int64                 `gorm:"column:update_time;not null;autoUpdateTime:milli" json:"update_time"`
	DeleteTime   soft_delete.DeletedAt `gorm:"column:delete_time;not null;softDelete:milli" json:"delete_time"`
}

func (UserInfo) TableName() string {
	return TableNameCommUserInfo
}
