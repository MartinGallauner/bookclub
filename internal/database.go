package internal

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	DbUser   string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Dbname   string `env:"POSTGRES_DBNAME"`
	Port     string `env:"POSTGRES_PORT"`
	Sslmode  string `env:"POSTGRES_SSLMODE"`
	Timezone string `env:"TIMEZONE"`
}

func ReadDatabaseConfig() (DatabaseConfig, error) {
	var DbConfig DatabaseConfig
	err := env.Parse(&DbConfig)
	if err != nil {
		return DatabaseConfig{}, err //consider returning pointer to struct
	}
	return DbConfig, nil
}

func SetupDatabase(config DatabaseConfig) (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", 
	config.Host, config.DbUser, config.Password, config.Dbname, config.Port, config.Sslmode, config.Timezone)
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
