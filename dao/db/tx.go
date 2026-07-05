package db

import (
	"context"
	"forum/dbcore"

	"gorm.io/gorm"
)

type txOpr struct {
	dbCtx dbcore.Context
}

func NewTxer(conn gorm.DB) Txer {
	return &txOpr{dbCtx: dbcore.NewContext(&conn)}
}

func (t *txOpr) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.dbCtx.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 把 *gorm.DB 绑定到 ctx，供 DBContext 使用
		return fn(dbcore.WithTx(ctx, tx))
	})
}
