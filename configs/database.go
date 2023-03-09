package configs

import (
	"errors"
	"fmt"
	"go-fiber-starter/utils"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func (c *Config) GormDatabase() {
	if os.Getenv("APP_DEBUG") == "" {
		err := errors.New("you should add APP_DEBUG variables, true or false")
		log.Fatal(err)
	}

	var (
		username string = os.Getenv("DB_USERNAME")
		password string = os.Getenv("DB_PASSWORD")
		url      string = os.Getenv("DB_HOST")
		port     string = os.Getenv("DB_PORT")
		dbName   string = os.Getenv("DB_NAME")
		debug    bool   = os.Getenv("DB_DEBUG") == "true"
	)

	logMode := logger.Silent
	if debug {
		logMode = logger.Info
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		url,
		port,
		username,
		dbName,
		password,
	)

	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		utils.Logger.Error("GORM Error : " + err.Error())
		log.Fatal("Open database connection failed")
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		log.Fatalf("Error connection database! %s", err.Error())
	}

	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxConnLifetime, _ := strconv.Atoi(os.Getenv("DB_MAX_CONN_LIFETIME"))
	maxConnIdleTime, _ := strconv.Atoi(os.Getenv("DB_MAX_CONN_IDLE_TIME"))

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxConnIdleTime) * time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Duration(maxConnLifetime) * time.Hour)

	DB = dbInstance
}
