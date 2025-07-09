package byjson

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"web-clean/infra/conf"
	"web-clean/infra/loader"
)

var JSONLoader loader.Loader = &_json{}

const (
	// MaxConfigFileSize 100 MB 是一个相当合理的限制的配置文件大小限制
	MaxConfigFileSize = 100 * 1024 * 1024
)

var MaxConfigFileSizeMB = toMB(MaxConfigFileSize)

func toMB(size int) int {
	return size / (1024 * 1024)
}

type _json struct {
}

func (_ _json) Load(ctx *loader.Context) (*conf.Conf, error) {
	return load(ctx)
}

func load(ctx *loader.Context) (*conf.Conf, error) {

	loadConfig := ctx.Config

	notEmptyFiles := filterEmptyString(loadConfig.Files)
	notEmptyPaths := filterEmptyString(loadConfig.Paths)

	if len(notEmptyPaths) == 0 || len(notEmptyFiles) == 0 {

		ctx.Log.Errorw("配置文件的路径或者文件名为空", "config", loadConfig)

		return nil, &Error{
			Msg: "配置文件路径或者文件名为空",
			Err: nil,
		}
	}

	errors := make([]error, 0)

	for _, p := range notEmptyPaths {
		for _, n := range notEmptyFiles {
			f := filepath.Join(p, n)

			ctx.Log.Debug("尝试配置文件路径")

			if _, err := os.Stat(f); os.IsNotExist(err) {

				ctx.Log.Debug("配置文件不存在")

				continue
			} else {
				ctx.Log.Debug("发现配置文件")

				parse, err := parse(ctx, f)

				if err != nil {
					errors = append(errors, err)
					continue
				} else {
					return parse, nil
				}
			}
		}
	}

	ctx.Log.Errorw("未找到可用的配置文件", "config", loadConfig, "errors", errors)

	// 如果存在 Parser 返回的错误那么返回最后一个，这是我们最好的选择
	if len(errors) != 0 {

		last := errors[len(errors)-1]

		return nil, &Error{
			Msg: last.Error(),
			Err: last,
		}
	}

	// 否则返回一个未找到
	return nil, &Error{
		Msg: "未找到可用的配置文件",
		Err: nil,
	}
}

func parse(ctx *loader.Context, path string) (*conf.Conf, error) {
	open, err := os.Open(path)
	if err != nil {
		ctx.Log.Errorw("无法打开文件", "error", err)
		return nil, &Error{
			Msg: "无法打开文件",
			Err: err,
		}
	}

	defer open.Close()

	fileStat, err := open.Stat()

	if err != nil {
		ctx.Log.Error("无法获得文件描述符", "error", err)
		return nil, &Error{
			Msg: "无法获得文件描述符",
			Err: err,
		}
	}

	if fileStat.Size() > MaxConfigFileSize {
		ctx.Log.Errorf("配置文件大小超过 %d", MaxConfigFileSizeMB)
		return nil, &Error{
			Msg: fmt.Sprintf("配置文件大小超过 %d", MaxConfigFileSizeMB),
			Err: nil,
		}
	}

	reader := io.LimitReader(open, MaxConfigFileSize)

	data, err := io.ReadAll(reader)

	if err != nil {
		ctx.Log.Error("无法读取配置文件", "error", err)
		return nil, &Error{
			Msg: "无法读取配置文件",
			Err: err,
		}
	}

	var config conf.Conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		ctx.Log.Error("无法反序列化配置文件到 Conf", "error", err)
		return nil, &Error{Msg: "无法反序列化配置文件到 Conf", Err: err}
	}

	return &config, nil
}

func filterEmptyString(strs []string) []string {
	notEmpty := make([]string, 0, len(strs))
	for _, str := range strs {
		if len(strings.TrimSpace(str)) != 0 {
			notEmpty = append(notEmpty, str)
		}
	}

	return notEmpty
}
