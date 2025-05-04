package cache

import (
	"gorm.io/driver/sqlite"
	"testing"
)

func TestCache(t *testing.T) {
	db := sqlite.Open("/Users/bytedance/.cache/golang-language-server/" + "go_lsp_cahce.db")
	//db, err := gorm.Open(sqlite.Open("/Users/bytedance/.cache/golang-language-server/"+conts.CacheFileName), &gorm.Config{
	//	SkipDefaultTransaction: false,
	//	PrepareStmt:            true,
	//})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	println(db)
}
