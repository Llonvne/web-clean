package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"web-clean/infra/web"
)

type Router struct {
	User
}

type Base struct {
	WebContextGetter func(ctx *gin.Context) (*web.Context, bool)
}

func (b *Base) WebContextMust(c *gin.Context) *web.Context {
	webCtx, ok := b.WebContextGetter(c)
	if !ok {
		// SAFETY: 常量错误永远不为 nil
		_ = c.Error(WebContextNotFoundError)

		// SAFETY: 该函数从不 return
		PanicServerInternalError(ServerInternalError{
			code: http.StatusInternalServerError,
			Body: gin.H{
				"message": "WebContextMust not found",
			},
			err: WebContextNotFoundError,
		})

		panic("unreachable")
	}

	return webCtx
}
