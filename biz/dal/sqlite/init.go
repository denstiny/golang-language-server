package sqlite

import (
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/denstiny/golang-language-server/biz/dal/sqlite/model"
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
	err = db.Table(model.PackageTableName).AutoMigrate(&model.Package{})
	if err != nil {
		panic("failed to migrate database" + err.Error())
	}
	err = db.Table(model.PackageLibranyTablName).AutoMigrate(&model.PackageLibrany{})
	if err != nil {
		panic("failed to migrate database" + err.Error())
	}
	err = db.Table(model.IndexTableName).AutoMigrate(&model.Index{})
	if err != nil {
		panic("failed to migrate database" + err.Error())
	}
}
