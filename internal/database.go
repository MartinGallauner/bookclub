package internal

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(host, dbUser, password, dbname, port, sslmode, timezone string) (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, dbUser, password, dbname, port, sslmode, timezone)
	return SetupDatabaseWithDSN(connString)
}

func SetupDatabaseWithDSN(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&User{}, &Book{}, &UserBooks{}, &Link{})
	if err != nil {
		return nil, err
	}
	err = db.SetupJoinTable(&User{}, "Books", &UserBooks{})
	if err != nil {
		return nil, err
	}
	var user User
	db.Preload("books").Find(&user)
	return db, nil
}
