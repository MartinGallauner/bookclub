package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host        string `env:"POSTGRES_HOST" validate:"hostname"`
	DbUser      string `env:"POSTGRES_USER"`
	Password    string `env:"POSTGRES_PASSWORD"`
	Dbname      string `env:"POSTGRES_DBNAME"`
	Port        int    `env:"POSTGRES_PORT"`
	Sslmode     string `env:"POSTGRES_SSLMODE"`
	Timezone    string `env:"TIMEZONE"`
	DatabaseUrl string `env:"DATABASE_URL"`
}

func ReadDatabaseConfig() (DatabaseConfig, error) {
	var dbConfig DatabaseConfig
	err := env.Parse(&dbConfig)
	if err != nil {
		return DatabaseConfig{}, err //consider returning pointer to struct
	}
	return dbConfig, nil
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func SetupDatabase(config DatabaseConfig) (*gorm.DB, error) {
	err := validate.Struct(&config)
	if err != nil {
		return nil, err
	}
	// commented out to use the connection string provided by Digital Ocean
	/* connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Host, config.DbUser, config.Password, config.Dbname, config.Port, config.Sslmode, config.Timezone) */
	return SetupDatabaseWithDSN(config.DatabaseUrl)
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
