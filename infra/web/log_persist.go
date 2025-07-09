package web

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"web-clean/domain"
)

type LogPersister interface {
	Persist(logs []Log) error
}

type Log struct {
	Level string
	Msg   string
}

type webLog struct {
	inner   domain.Log
	context *gin.Context
	logs    []Log
}

func (w *webLog) appendToLogs(level string, args ...interface{}) {
	// 将 args 序列化为 JSON
	jsonBytes, err := json.Marshal(args)
	var msg string
	if err != nil {
		msg = fmt.Sprintf("序列化错误: %v, 原始参数: %v", err, args)
	} else {
		msg = string(jsonBytes)
	}

	// 添加到日志切片
	w.logs = append(w.logs, Log{
		Level: level,
		Msg:   msg,
	})
}

func (w *webLog) Debug(args ...interface{}) {
	w.inner.Debug(args...)
	w.appendToLogs("DEBUG", args...)
}

func (w *webLog) Info(args ...interface{}) {
	w.inner.Info(args...)
	w.appendToLogs("INFO", args...)
}

func (w *webLog) Warn(args ...interface{}) {
	w.inner.Warn(args...)
	w.appendToLogs("WARN", args...)
}

func (w *webLog) Error(args ...interface{}) {
	w.inner.Error(args...)
	w.appendToLogs("ERROR", args...)
}

func (w *webLog) DPanic(args ...interface{}) {
	w.inner.DPanic(args...)
	w.appendToLogs("DPANIC", args...)
}

func (w *webLog) Panic(args ...interface{}) {
	w.inner.Panic(args...)
	w.appendToLogs("PANIC", args...)
}

func (w *webLog) Fatal(args ...interface{}) {
	w.inner.Fatal(args...)
	w.appendToLogs("FATAL", args...)
}

func (w *webLog) Debugf(template string, args ...interface{}) {
	w.inner.Debugf(template, args...)
	w.appendToLogs("DEBUG", template, args)
}

func (w *webLog) Infof(template string, args ...interface{}) {
	w.inner.Infof(template, args...)
	w.appendToLogs("INFO", template, args)
}

func (w *webLog) Warnf(template string, args ...interface{}) {
	w.inner.Warnf(template, args...)
	w.appendToLogs("WARN", template, args)
}

func (w *webLog) Errorf(template string, args ...interface{}) {
	w.inner.Errorf(template, args...)
	w.appendToLogs("ERROR", template, args)
}

func (w *webLog) DPanicf(template string, args ...interface{}) {
	w.inner.DPanicf(template, args...)
	w.appendToLogs("DPANIC", template, args)
}

func (w *webLog) Panicf(template string, args ...interface{}) {
	w.inner.Panicf(template, args...)
	w.appendToLogs("PANIC", template, args)
}

func (w *webLog) Fatalf(template string, args ...interface{}) {
	w.inner.Fatalf(template, args...)
	w.appendToLogs("FATAL", template, args)
}

func (w *webLog) Debugw(msg string, keysAndValues ...interface{}) {
	w.inner.Debugw(msg, keysAndValues...)
	w.appendToLogs("DEBUG", msg, keysAndValues)
}

func (w *webLog) Infow(msg string, keysAndValues ...interface{}) {
	w.inner.Infow(msg, keysAndValues...)
	w.appendToLogs("INFO", msg, keysAndValues)
}

func (w *webLog) Warnw(msg string, keysAndValues ...interface{}) {
	w.inner.Warnw(msg, keysAndValues...)
	w.appendToLogs("WARN", msg, keysAndValues)
}

func (w *webLog) Errorw(msg string, keysAndValues ...interface{}) {
	w.inner.Errorw(msg, keysAndValues...)
	w.appendToLogs("ERROR", msg, keysAndValues)
}

func (w *webLog) DPanicw(msg string, keysAndValues ...interface{}) {
	w.inner.DPanicw(msg, keysAndValues...)
	w.appendToLogs("DPANIC", msg, keysAndValues)
}

func (w *webLog) Panicw(msg string, keysAndValues ...interface{}) {
	w.inner.Panicw(msg, keysAndValues...)
	w.appendToLogs("PANIC", msg, keysAndValues)
}

func (w *webLog) Fatalw(msg string, keysAndValues ...interface{}) {
	w.inner.Fatalw(msg, keysAndValues...)
	w.appendToLogs("FATAL", msg, keysAndValues)
}
