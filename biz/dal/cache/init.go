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
	db, err := gorm.Open(sqlite.Open(path.Join(flags.SERVICE_CONFIG_DIR, conts.CacheFileName)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}
