package model

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/dynamicmapping"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	TableNamePost = "post"
)

// Post 社区内容统一模型（帖子与评论融合，纯 parent_id 邻接树，最多 2 层）
//   - parent_id = 0          : 帖子（根节点）
//   - parent_id = 帖子id     : 一级评论
//   - parent_id = 一级评论id : 回复
type Post struct {
	ID               int64                       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID         int64                       `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	UserID           int64                       `gorm:"column:user_id;not null;default:0" json:"user_id"`
	Account          string                      `gorm:"column:account;type:varchar(100);not null;default:''" json:"account"`
	Content          string                      `gorm:"column:content;type:text;not null" json:"content"`
	Images           datatypes.JSONSlice[string] `gorm:"column:images;type:json;not null;default:[]" json:"images"`
	WordCount        int32                       `gorm:"column:word_count;not null;default:0" json:"word_count"`
	Status           int32                       `gorm:"column:status;not null;default:30" json:"status"` // 30:已发布 90:已删除
	Remark           *string                     `gorm:"column:remark;type:varchar(255);not null;default:''" json:"remark"`
	ViewCount        int32                       `gorm:"column:view_count;not null;default:0" json:"view_count"`
	LikeCount        int32                       `gorm:"column:like_count;not null;default:0" json:"like_count"`
	CommentCount     int32                       `gorm:"column:comment_count;not null;default:0" json:"comment_count"`
	FavCount         int32                       `gorm:"column:fav_count;not null;default:0" json:"fav_count"`
	IPAddress        string                      `gorm:"column:ip_address;type:varchar(256);not null;default:''" json:"ip_address"`
	PublishTime      *int64                      `gorm:"column:publish_time;not null;default:0" json:"publish_time"`
	StatusUpdateTime *int64                      `gorm:"column:status_update_time;not null;default:0" json:"status_update_time"`
	CreateTime       int64                       `gorm:"column:create_time;not null;autoCreateTime:milli" json:"create_time"`
	UpdateTime       int64                       `gorm:"column:update_time;not null;autoUpdateTime:milli" json:"update_time"`
	DeleteTime       soft_delete.DeletedAt       `gorm:"column:delete_time;not null;softDelete:milli" json:"delete_time"`
	// Parent 直接父节点（评论的父帖子或父评论）；不参与 ES 序列化
	Parent *Post `gorm:"foreignKey:ParentID;references:ID" json:"-"`
}

func (Post) TableName() string {
	return TableNamePost
}

// BeforeSave 确保 images 为空时写入空数组("[]")而不是 JSON null
func (p *Post) BeforeSave(*gorm.DB) error {
	if p.Images == nil {
		p.Images = datatypes.NewJSONSlice([]string{})
	}
	return nil
}

func (Post) Mapping() *types.TypeMapping {
	return &types.TypeMapping{
		Dynamic: &dynamicmapping.Strict,
		Properties: map[string]types.Property{
			"id":                 types.NewLongNumberProperty(),
			"parent_id":          types.NewLongNumberProperty(),
			"user_id":            types.NewLongNumberProperty(),
			"account":            types.NewKeywordProperty(),
			"content":            types.NewTextProperty(),
			"images":             types.NewKeywordProperty(),
			"word_count":         types.NewIntegerNumberProperty(),
			"status":             types.NewIntegerNumberProperty(),
			"remark":             types.NewTextProperty(),
			"view_count":         types.NewIntegerNumberProperty(),
			"like_count":         types.NewIntegerNumberProperty(),
			"comment_count":      types.NewIntegerNumberProperty(),
			"fav_count":          types.NewIntegerNumberProperty(),
			"ip_address":         types.NewKeywordProperty(),
			"publish_time":       types.NewLongNumberProperty(),
			"status_update_time": types.NewLongNumberProperty(),
			"create_time":        types.NewLongNumberProperty(),
			"update_time":        types.NewLongNumberProperty(),
			"delete_time":        types.NewLongNumberProperty(),
		},
	}
}
