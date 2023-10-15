package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	//	subject to change
	return db.AutoMigrate(&User{})
}
