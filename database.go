package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&User{}, &Book{}, &UserBooks{})
	err = db.SetupJoinTable(&User{}, "Books", &UserBooks{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
