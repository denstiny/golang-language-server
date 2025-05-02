package model

import "strings"

const PackageTableName = "packages"
const IndexTableName = "indexes"

/*
Package 结构体用于表示一个软件包的相关信息，这些信息会被存储到数据库中，并且能够以 JSON 格式进行输出。
结构体的各个字段分别从不同方面描述了软件包的特征，具体说明如下：

  - Name: 软件包的名称，是软件包的重要标识之一。该字段在数据库中对应的字段名为 "name"，在进行 JSON 序列化时键名为 "name"。
    其数据库类型为 varchar(1024)，并且为该字段创建了名为 idx_name 的索引，以加快基于包名的查询速度。
  - Repo: 软件包所在的仓库地址，指明了软件包的来源位置。在数据库中对应的字段名为 "repo"，JSON 序列化时键名为 "repo"。
    数据库类型为 varchar(1024)，同时创建了名为 idx_repo 的索引，方便对仓库地址相关的查询操作。
  - Version: 软件包的版本号，用于区分同一软件包的不同迭代版本。数据库字段名为 "version"，JSON 键名为 "version"。
    数据库类型为 varchar(1024)，并创建了名为 idx_version 的索引，有助于提高基于版本号的查询效率。
*/
type Package struct {
	ID      int64  `db:"id" json:"id" gorm:"primary_key"`
	Name    string `db:"name" json:"name" gorm:"type:varchar(1024)index:idx_name"`
	Repo    string `db:"repo" json:"repo" gorm:"type:varchar(1024)index:idx_repo"`
	Version string `db:"version" json:"version" gorm:"type:varchar(1024)index:idx_version"`
}

/*
Index 结构体代表了索引的相关信息，它会存储在数据库里，同时也能以 JSON 格式输出。
该结构体包含了以下字段，用于描述索引的详细信息：
  - Comparable: 可用于比较的内容，通常是索引中某个关键的可对比元素，在数据库中对应 "comparable" 字段，JSON 序列化时键名为 "comparable"。
  - KeyWorld: 索引的关键字，是索引的核心标识之一，在数据库中对应 "key_world" 字段，JSON 序列化时键名为 "key_world"。
  - Type: 索引的类型，表明索引所属的类别，在数据库中对应 "type" 字段，JSON 序列化时键名为 "type"。
  - JoinIndex: 存储 Join 操作跳转的文件位置，为数据关联操作提供跳转依据，在数据库中对应 "join_index" 字段，JSON 序列化时键名为 "join_index"。
  - FilePath: 索引所在的文件位置，指明索引信息来源于哪个文件，在数据库中对应 "file_path" 字段，JSON 序列化时键名为 "file_path"。
  - Package: 索引所在的包名，明确索引所属的 Go 包，在数据库中对应 "package" 字段，JSON 序列化时键名为 "package"。
  - JoinLine: 索引所在的行号，精确到文件中的行位置，在数据库中对应 "join_line" 字段，JSON 序列化时键名为 "join_line"。
  - JoinCol: 索引所在的列号，精确到文件中的列位置，在数据库中对应 "join_col" 字段，JSON 序列化时键名为 "join_col"。
*/
type Index struct {
	ID         int    `db:"id" json:"id" gorm:"primary_key"`
	Comparable string `db:"comparable" json:"comparable" gorm:"type:text"`
	KeyWorld   string `db:"key_world" json:"key_world" gorm:"type:varchar(1024)index:idx_key_world"`
	Type       int32  `db:"type" json:"type" gorm:"type:int:index:idx_type"`
	JoinIndex  string `db:"join_index" json:"join_index"`
	FilePath   string `db:"file_path" json:"file_path" gorm:"type:varchar(2048)index:idx_file_path"`
	Package    string `db:"package" json:"package" gorm:"type:varchar(1024)index:idx_package"`
	JoinLine   int    `db:"join_line" json:"join_line" gorm:"type:int"`
	JoinCol    int    `db:"join_col" json:"join_col" gorm:"type:"`
	PackageID  int32  `db:"package_id" json:"package_id" gorm:"type:int index:idx_package_id"`
}

func (p *Package) IndexName() string {
	return strings.Join([]string{p.Name, p.Version}, "@")
}
