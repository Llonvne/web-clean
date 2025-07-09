package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"web-clean/infra"
	"web-clean/infra/database"
	"web-clean/infra/web"
)

type ErrorModel struct {
	gorm.Model
	Error web.Errors `gorm:"type:jsonb"`
}

func init() {
	database.RegisterSchema(ErrorModel{})
}

type Errors struct {
	*infra.Context

	FallbackFilePath string

	Database database.Database
}

func (e Errors) Persist(errors web.Errors) {
	err := e.Database.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ErrorModel{
			Error: errors,
		}).Error
	})
	if err != nil {
		err := e.saveToFile(errors)

		if err != nil {
			e.Log.Errorw("无法向错误数据库写入错误堆栈，也无法向错误文件写入", "err", err, "errors", errors)
		}
	}
}

func (e Errors) saveToFile(rec web.Errors) error {
	data, _ := json.MarshalIndent(rec, "", "  ")

	if err := os.MkdirAll(e.FallbackFilePath, 0o755); err != nil {
		return err
	}

	fileName := fmt.Sprintf("error_%s_%s.json",
		rec.RequestID,
		time.Now().Format("20060102T150405.000"),
	)

	fullPath := filepath.Join(e.FallbackFilePath, fileName)

	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}
