package db

import (
	"forum/types/forum"

	"gorm.io/gorm"
)

type ScopeF = func(db *gorm.DB) *gorm.DB

func OffsetAndLimitScope(offset, limit int) ScopeF {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

func SortScope(sorts []*forum.Sort) ScopeF {
	return func(db *gorm.DB) *gorm.DB {
		if len(sorts) == 0 {
			return db.Order("id DESC")
		}
		for _, sort := range sorts {
			db = db.Order(sort.Field + " " + sort.Order)
		}
		return db
	}
}

// ParentIdEqScope 过滤指定父节点的直接子节点（parent_id = 0 即帖子）
func ParentIdEqScope(parentId int64) ScopeF {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("parent_id = ?", parentId)
	}
}

// ParentIdNotZeroScope 仅保留评论（parent_id <> 0）
func ParentIdNotZeroScope() ScopeF {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("parent_id <> 0")
	}
}

// TargetIdsInScope 过滤点赞/收藏的目标内容集合
func TargetIdsInScope(targetIds []int64) ScopeF {
	return func(db *gorm.DB) *gorm.DB {
		if len(targetIds) == 0 {
			return db
		}
		return db.Where("target_id IN ?", targetIds)
	}
}

// PreloadTargetScope 预加载收藏/点赞的目标内容
func PreloadTargetScope() ScopeF {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Target")
	}
}
