package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	db   *gorm.DB
	err  error
	once sync.Once
)

const (
	DSN = "root:123456@tcp(127.0.0.1:3306)/money?charset=utf8mb4&parseTime=True&loc=Local"
)

func DB() *gorm.DB {
	once.Do(func() {
		db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	})
	return db
}
