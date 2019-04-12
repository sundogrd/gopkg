package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // ...
)

type ConnectOptions struct {
	User string
	Password string
	Host string
	Port string
	DBName string
	ConnectTimeout string
}

// OpenDB 用连接字符串连接数据库
func OpenDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Connect 用配置文件连接数据库
func Connect(connectOptions ConnectOptions) (*gorm.DB, error) {
	db, err := OpenDB(connectOptions.User + ":" + connectOptions.Password + "@tcp(" + connectOptions.Host + ":" + connectOptions.Port + ")" + "/" + connectOptions.DBName + "?charset=utf8&parseTime=True&loc=Local&timeout=" + connectOptions.ConnectTimeout)
	if err != nil {
		return nil, err
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	db.LogMode(true)
	return db, nil
}