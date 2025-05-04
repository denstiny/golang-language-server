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
	Version     *string
	Id          *int32
	PackageName *string
	Name        *string
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
	if find.PackageName != nil {
		db = db.Where("package_name=?", *find.PackageName)
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

func GetPackage(packageName string, version string) (*model.Package, error) {
	db := DB.Table(model.PackageTableName)
	db = db.Where("package_name=? and version=?", packageName, version)
	var p model.Package
	err := db.First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func FindPackageLibrany(packageId int64) ([]*model.PackageLibrany, error) {
	db := DB.Table(model.PackageLibranyTablName)
	err := db.AutoMigrate()
	if err != nil {
		return nil, err
	}
	db.Where("package_id=?", packageId)
	var results []*model.PackageLibrany
	err = db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil

}

func CreatePackageLibrany(pg model.PackageLibrany) error {
	db := DB.Table(model.PackageLibranyTablName)
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return db.Create(&pg).Error
}

func FindPackageLibranyLikeName(name string) ([]*model.PackageLibrany, error) {
	db := DB.Table(model.PackageLibranyTablName)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	for _, c := range name {
		db = db.Where("name like ?", "%"+string(c)+"%")
	}

	var results []*model.PackageLibrany
	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
