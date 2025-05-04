package cache

import (
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/denstiny/golang-language-server/biz/flags"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
)

var (
	DB *gorm.DB
)

func init() {
	dbpath := path.Join(flags.SERVICE_CONFIG_DIR, conts.CacheFileName)
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
