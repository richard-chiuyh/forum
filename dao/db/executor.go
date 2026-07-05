package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CountColumn string

const (
	COMMENT_COUNT CountColumn = "comment_count"
	FAV_COUNT     CountColumn = "fav_count"
	LIKE_COUNT    CountColumn = "like_count"
)

type CountTask struct {
	Table  string      // 目标表名
	Id     int64       // 记录 ID
	Column CountColumn // 要累加的列名
	Delta  int64       // 增量
}

type Executor struct {
	executor *executors.BulkExecutor
}

func (e *Executor) Add(task CountTask) {
	e.executor.Add(task)
}

func (e *Executor) Flush() {
	e.executor.Flush()
}

func NewExecutor(conn gorm.DB) Executor {
	executor := executors.NewBulkExecutor(func(tasks []interface{}) {
		if len(tasks) == 0 {
			return
		}
		group := make(map[string]map[CountColumn]map[int64]int64)
		for _, t := range tasks {
			task := t.(CountTask)
			if group[task.Table] == nil {
				group[task.Table] = make(map[CountColumn]map[int64]int64)
			}
			if group[task.Table][task.Column] == nil {
				group[task.Table][task.Column] = make(map[int64]int64)
			}
			group[task.Table][task.Column][task.Id] += task.Delta
		}

		db := conn.Session(&gorm.Session{})
		for table, columnMap := range group {
			ids := collectIDs(columnMap)
			updates := make(map[string]interface{}, len(columnMap))
			for col, idMap := range columnMap {
				updates[string(col)] = gorm.Expr(buildCaseExpr(col, idMap))
			}
			if err := db.Table(table).Where("id IN ?", ids).UpdateColumns(updates).Error; err != nil {
				logx.Errorf("CountExecutor table(%s) ids(%v) err(%+v)", table, ids, err)
			}
		}
	}, executors.WithBulkInterval(time.Second*1), executors.WithBulkTasks(128))
	return Executor{executor: executor}
}

func collectIDs(columnMap map[CountColumn]map[int64]int64) []int64 {
	set := make(map[int64]struct{})
	for _, idMap := range columnMap {
		for id := range idMap {
			set[id] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}

func buildCaseExpr(col CountColumn, idMap map[int64]int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "`%s` + CASE `id` ", col)
	for id, delta := range idMap {
		fmt.Fprintf(&b, "WHEN %d THEN %d ", id, delta)
	}
	b.WriteString("ELSE 0 END")
	return b.String()
}
