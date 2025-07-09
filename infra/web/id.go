package web

import "github.com/gin-gonic/gin"

var (
	idKey = "__idKey__"
)

func RequestIDMiddleware(idGen func() string) gin.HandlerFunc {
	return func(context *gin.Context) {
		requestID := context.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = idGen()
			context.Set("X-Request-ID", requestID)
		}
		context.Set(idKey, requestID)
		context.Next()
	}
}

func RequestIdGetter(c *gin.Context) string {
	s, ok := c.Get(idKey)
	if !ok {
		return ""
	}
	return s.(string)
}
