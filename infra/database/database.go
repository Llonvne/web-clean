package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"web-clean/infra"
)

type Database interface {
	Transaction(func(tx *gorm.DB) error) error
}

type _database struct {
	raw *gorm.DB
}

func (d *_database) Transaction(f func(tx *gorm.DB) error) error {
	return d.raw.Transaction(f)
}

func From(ctx *infra.Context) (Database, error) {

	config := ctx.Conf.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host, config.Username, config.Password, config.Database, config.Port)

	ctx.Log.Infow("连接PostgresSQL数据库", "dsn", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &_database{raw: db}, nil
}
