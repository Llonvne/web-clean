package web

import "github.com/gin-gonic/gin"

func Recover(panicHandler func(context *gin.Context, err any)) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				panicHandler(context, err)
			}
		}()
		context.Next()
	}
}
