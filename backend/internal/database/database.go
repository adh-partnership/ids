package database

import (
	"database/sql"
	"fmt"
	"time"

	"dario.cat/mergo"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type DBOptions struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  string

	MaxOpenConns int
	MaxIdleConns int

	CACert string
}

var defaultOptions = DBOptions{
	MaxOpenConns: 50,
	MaxIdleConns: 10,
}

func GenerateDSN(options DBOptions) (string, error) {
	var dsn string

	if options.Driver == "mysql" {
		tls := ""
		if options.CACert != "" {
			tls = "&tls=custom"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s", options.User, options.Password,
			options.Host, options.Port, options.Database, tls)
		if options.Options != "" {
			dsn += "?" + options.Options
		}
	} else if options.Driver == "postgres" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=prefer",
			options.Host, options.Port, options.User, options.Password, options.Database)
	} else {
		return "", fmt.Errorf("unsupported driver: %s", options.Driver)
	}

	return dsn, nil
}

func IsValidDriver(driver string) bool {
	return driver == "mysql" || driver == "postgres"
}

func New(options DBOptions) (*Database, error) {
	if !IsValidDriver(options.Driver) {
		return nil, fmt.Errorf("invalid driver: %s", options.Driver)
	}

	err := mergo.Merge(&options, defaultOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to merge options: %w", err)
	}

	dsn, err := GenerateDSN(options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dsn: %w", err)
	}

	var db *gorm.DB
	config := &gorm.Config{
		Logger: NewGormLogger(),
	}
	if options.Driver == "mysql" {
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open mysql connection: %w", err)
		}

		db, err = gorm.Open(mysql.New(mysql.Config{Conn: conn}), config)
		if err != nil {
			return nil, fmt.Errorf("failed to open gorm connection: %w", err)
		}
	} else if options.Driver == "postgres" {
		db, err = gorm.Open(postgres.Open(dsn), config)
		if err != nil {
			return nil, fmt.Errorf("failed to open gorm connection: %w", err)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(options.MaxIdleConns)
	sqlDB.SetMaxOpenConns(options.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return &Database{DB: db}, nil
}
