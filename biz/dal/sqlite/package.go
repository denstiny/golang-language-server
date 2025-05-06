package sqlite

import "github.com/denstiny/golang-language-server/biz/dal/sqlite/model"

func CreatePackage(pg *model.Package) error {
	db := DB.Table(model.PackageTableName)
	return db.Create(pg).Error
}

type PackageFindParams struct {
	Version     *string
	Id          *int32
	PackagePath *string
	Name        *string
}

func FindPackage(find PackageFindParams) ([]*model.Package, error) {
	db := DB.Table(model.PackageTableName)
	if find.Id != nil {
		db = db.Where("id=?", *find.Id)
	}
	if find.Version != nil {
		db = db.Where("version=?", *find.Version)
	}
	if find.PackagePath != nil {
		db = db.Where("package_path=?", *find.PackagePath)
	}
	if find.Name != nil {
		db = db.Where("name=?", *find.Name)
	}
	var results []*model.Package
	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetPackage(packagePath string, version string) (*model.Package, error) {
	db := DB.Table(model.PackageTableName)
	db = db.Where("package_path=? and version=?", packagePath, version)
	var p model.Package
	err := db.First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func FindPackageLibrany(packageId int64) ([]*model.PackageLibrany, error) {
	db := DB.Table(model.PackageLibranyTablName)
	db.Where("package_id=?", packageId)
	var results []*model.PackageLibrany
	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil

}

func CreatePackageLibrany(pg *model.PackageLibrany) error {
	db := DB.Table(model.PackageLibranyTablName)
	return db.Create(pg).Error
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
