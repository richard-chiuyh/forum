package es

import (
	"context"

	"forum/dao/model"
	"forum/escore"
	"forum/internal/config"
	"forum/types/forum"

	"github.com/zeromicro/go-zero/core/logx"
)

type Manager struct {
	PostDao PostDao
}

func NewManager(conf config.Config) (Manager, error) {
	esClient, err := escore.Connect(&conf.ES)
	if err != nil {
		logx.Severef("es client init err(%+v)", err)
		return Manager{}, err
	}
	if err := initPostIndex(esClient, conf.Mode); err != nil {
		logx.Severef("es init index err(%+v)", err)
		return Manager{}, err
	}
	return Manager{
		PostDao: newPostOpr(esClient, conf.Mode),
	}, nil
}

type PostDao interface {
	Index(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, in *forum.ListPostRequest) ([]*model.Post, int64, error)
}
