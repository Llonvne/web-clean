package web

import (
	"github.com/gin-gonic/gin"

	"web-clean/domain"
	"web-clean/infra/database"
)

type Context struct {
	database.Database
	Log domain.Log
}

var (
	webContextKey = "__webCtxKey__"
)

func ContextMiddleware(
	constructor func(log domain.Log) *Context,
	innerLogger domain.Log,
	webLogPersister LogPersister,
) gin.HandlerFunc {

	return func(context *gin.Context) {

		webLogger := webLog{
			inner:   innerLogger,
			context: context,
			logs:    make([]Log, 0),
		}

		defer func() {
			webLogPersister.Persist(webLogger.logs)
		}()

		webCtx := constructor(&webLogger)

		context.Set(webContextKey, webCtx)
		defer context.Set(webContextKey, nil)

		context.Next()
	}
}

func ContextMiddlewareGetter(context *gin.Context) (*Context, bool) {
	value, exists := context.Get(webContextKey)
	if !exists {
		return nil, false
	}

	c, ok := value.(*Context)
	if !ok {
		return nil, false
	}

	return c, true
}
