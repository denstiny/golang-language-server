package cache

import (
	"github.com/denstiny/golang-language-server/biz/dal/cache/model"
)

func QueryIndexByPackageID(PackageID int) ([]*model.Index, error) {
	db := DB.Table(model.IndexTableName)
	err := db.AutoMigrate()
	if err != nil {
		return nil, err
	}

	var results []*model.Index
	err = db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

type IndexFindParams struct {
	Keyword   *string
	Filename  *string
	PackageID *int32
	Type      *int32
}

func FindIndex(params IndexFindParams) ([]*model.Index, error) {
	db := DB.Table(model.IndexTableName)
	err := db.AutoMigrate()
	if err != nil {
		return nil, err
	}
	if params.PackageID != nil {
		db = db.Where("package_id = ?", params.PackageID)
	}

	if params.Filename != nil {
		db = db.Where("filename = ?", *params.Filename)
	}
	if params.Type != nil {
		db = db.Where("type = ?", *params.Type)
	}

	if params.Keyword != nil {
		db = db.Where("keyword = ?", *params.Keyword)
	}
	var results []*model.Index
	err = db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func CreateIndex(index model.Index) error {
	db := DB.Table(model.IndexTableName)
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return db.Create(&index).Error
}
