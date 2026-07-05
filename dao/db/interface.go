package db

import (
	"context"

	"forum/dao/model"
	"forum/dbcore"

	"github.com/zeromicro/go-zero/core/logx"
)

type Manager struct {
	PostDao     PostDao
	UserInfoDao UserInfoDao
	FavDao      FavDao
	LikeDao     LikeDao
	Executor    Executor
	Txer        Txer
}

func NewManager(conf *dbcore.MysqlConfig) (Manager, error) {
	conn, err := dbcore.Connect(conf)
	if err != nil {
		logx.Severef("db connection err(%+v)", err)
		return Manager{}, err
	}
	sqlDB, err := conn.DB()
	if err != nil {
		logx.Severef("db connection err(%+v)", err)
		return Manager{}, err
	}
	if sqlDB.Ping() != nil {
		logx.Errorf("db ping err(%+v)", err)
		return Manager{}, err
	}
	return Manager{
		PostDao:     newPostOpr(*conn),
		UserInfoDao: newUserInfoOpr(*conn),
		FavDao:      newFavOpr(*conn),
		LikeDao:     newLikeOpr(*conn),
		Executor:    NewExecutor(*conn),
		Txer:        NewTxer(*conn),
	}, nil
}

// 内容操作接口（post 单表，帖子 parent_id=0 / 评论 parent_id<>0 统一走此接口）
type PostDao interface {
	Create(ctx context.Context, post *model.Post) error
	Update(ctx context.Context, post *model.Post) error
	Get(ctx context.Context, id int64) (*model.Post, error)
	// cond 用于简单等值条件（零值字段自动跳过），复杂条件通过 scopes 传入
	List(ctx context.Context, cond *model.Post, scopes ...ScopeF) ([]*model.Post, int64, error)
}

// 用户信息操作接口
type UserInfoDao interface {
	Create(ctx context.Context, userInfo *model.UserInfo) error
	Update(ctx context.Context, userInfo *model.UserInfo) error
	GetByUserID(ctx context.Context, userID int64) (*model.UserInfo, error)
}

// 收藏操作接口（target_id 指向内容）
type FavDao interface {
	Create(ctx context.Context, fav *model.Fav) error
	Get(ctx context.Context, userId, targetId int64) (*model.Fav, error)
	List(ctx context.Context, cond *model.Fav, scopes ...ScopeF) ([]*model.Fav, int64, error)
	Update(ctx context.Context, fav *model.Fav) error
}

// 点赞操作接口（target_id 指向内容）
type LikeDao interface {
	Create(ctx context.Context, like *model.Like) error
	Get(ctx context.Context, userId, targetId int64) (*model.Like, error)
	List(ctx context.Context, cond *model.Like, scopes ...ScopeF) ([]*model.Like, int64, error)
	Update(ctx context.Context, like *model.Like) error
}

type Txer interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
