package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"web-clean/infra"
)

type _gin struct {
	*infra.Context
	engine *gin.Engine
}

func (g *_gin) Serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", g.Conf.Web.Port),
		Handler: g.engine,
	}

	ctx, stop := signal.NotifyContext(g.Ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			g.Log.Fatal("💥 listen: ", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	stop()
	g.Log.Info("🛑 shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		g.Log.Fatal("💥 Server forced to shutdown: ", err)
	}

	g.Log.Infow("👋 Server exiting")
}

func Gin(
	ctx *infra.Context,
	opt func(*gin.Engine),
) Web {
	var g = &_gin{
		engine:  gin.New(),
		Context: ctx,
	}

	opt(g.engine)

	return g
}
