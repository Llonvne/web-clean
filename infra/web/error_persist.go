package web

import (
	"github.com/gin-gonic/gin"

	"web-clean/domain"
)

type Errors struct {
	Stack     any
	Method    string
	URL       string
	Path      string
	IP        string
	RequestID string
}

type ErrorStackPersister interface {

	// Persist 该方法将错误堆栈持久化，请注意由于该插件注册于 Recover 外围，该插件不能返回错误，不能 panic
	Persist(errors Errors)
}

func ErrorPersisterMiddleware(
	persistent ErrorStackPersister,
	log domain.Log,
	requestIdGetter func(ctx *gin.Context) string,
) gin.HandlerFunc {
	return func(context *gin.Context) {

		requestMethod := context.Request.Method    // 请求方法
		requestURL := context.Request.URL.String() // 完整 URL
		requestPath := context.Request.URL.Path    // 请求路径
		requestIP := context.ClientIP()            // 客户端 IP
		requestID := requestIdGetter(context)

		defer func() {
			// 理论上 persist 不应该出现错误，但是如果他 panic 我们能做的就是简单将错误和错误堆栈打印日志
			// 这已经是尽最大努力了
			err := recover()
			if err != nil {
				log.Errorw("ErrorPersister 抛出错误", "error", err, "stackErrors", context.Errors)
			}
		}()

		context.Next()

		if len(context.Errors) == 0 {
			return
		}

		var errBody any

		errBody = context.Errors.String()

		persistent.Persist(Errors{
			Stack:     errBody,
			Method:    requestMethod,
			URL:       requestURL,
			Path:      requestPath,
			IP:        requestIP,
			RequestID: requestID,
		})
	}
}
