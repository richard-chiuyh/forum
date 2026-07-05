package model

import (
	"gorm.io/plugin/soft_delete"
)

const TableNameLike = "like"

type Like struct {
	ID         int64                 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TargetID   int64                 `gorm:"column:target_id;not null;default:0" json:"target_id"` // 目标内容ID（帖子或评论）
	UserID     int64                 `gorm:"column:user_id;not null;default:0" json:"user_id"`
	Account    string                `gorm:"column:account;type:varchar(32);not null;default:''" json:"account"`
	Status     int32                 `gorm:"column:status;not null;default:0" json:"status"` // 30:正常 90:已删除
	IPAddress  string                `gorm:"column:ip_address;type:varchar(256);not null;default:''" json:"ip_address"`
	CreateTime int64                 `gorm:"column:create_time;not null;autoCreateTime:milli" json:"create_time"`
	UpdateTime int64                 `gorm:"column:update_time;not null;autoUpdateTime:milli" json:"update_time"`
	DeleteTime soft_delete.DeletedAt `gorm:"column:delete_time;not null;softDelete:milli" json:"delete_time"`
}

func (Like) TableName() string {
	return TableNameLike
}
