package dbcore

import (
	"context"

	"gorm.io/gorm"
)

type Context struct {
	conn *gorm.DB
}

func NewContext(conn *gorm.DB) Context {
	return Context{conn: conn}
}

type txContextKey struct{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func TxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

func (p *Context) DB(ctx context.Context) *gorm.DB {
	if tx := TxFromContext(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return p.conn.Session(&gorm.Session{NewDB: true}).WithContext(ctx)
}
