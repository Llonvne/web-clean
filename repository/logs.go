package repository

import (
	"gorm.io/gorm"

	"web-clean/infra"
	"web-clean/infra/database"
	"web-clean/infra/web"
)

type LogsModel struct {
	gorm.Model
	Logs []web.Log `gorm:"type:jsonb"`
}

func init() {
	database.RegisterSchema(LogsModel{})
}

type Logs struct {
	*infra.Context
	database.Database
}

func (l *Logs) Persist(logs []web.Log) error {
	return l.Database.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&LogsModel{
			Logs: logs,
		}).Error
	})
}
