package gfadapter

import "github.com/gogf/gf/database/gdb"

func NewAdapterByGdb(db gdb.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}
