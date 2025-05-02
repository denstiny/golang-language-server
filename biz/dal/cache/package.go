package cache

import "github.com/denstiny/golang-language-server/biz/dal/cache/model"

func CreatePackage(pg model.Package) error {
	db := DB.Table(model.PackageTableName)
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return db.Create(&pg).Error
}

type PackageFindParams struct {
	Version *string
	Id      *int32
	Repo    *string
	Name    *string
}

func FindPackage(find PackageFindParams) ([]*model.Package, error) {
	db := DB.Table(model.PackageTableName)
	err := db.AutoMigrate()
	if err != nil {
		return nil, err
	}
	if find.Id != nil {
		db = db.Where("id=?", *find.Id)
	}
	if find.Version != nil {
		db = db.Where("version=?", *find.Version)
	}
	if find.Repo != nil {
		db = db.Where("repo=?", *find.Repo)
	}
	if find.Name != nil {
		db = db.Where("name=?", *find.Name)
	}
	var results []*model.Package
	err = db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
