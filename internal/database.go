package internal

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&User{}, &Book{}, &UserBooks{}, &Link{})
	err = db.SetupJoinTable(&User{}, "Books", &UserBooks{})
	if err != nil {
		return nil, err
	}
	var user User
	db.Preload("books").Find(&user)
	return db, nil
}
