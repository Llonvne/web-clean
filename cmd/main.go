package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"web-clean/domain"
	"web-clean/handler"
	"web-clean/infra"
	"web-clean/infra/database"
	byjson "web-clean/infra/loader/json"
	"web-clean/infra/web"
	"web-clean/repository"
)

func main() {
	context, err := infra.Prepare(infra.PrepareConfig{Loader: byjson.JSONLoader})
	if err != nil {
		panic(err)
	}

	db, err := database.From(context)
	if err != nil {
		panic(err)
	}

	err = database.AutoMigrateRegisteredSchema(db)
	if err != nil {
		panic(err)
	}

	logsPersister := repository.Logs{
		Context:  context,
		Database: db,
	}

	contextMiddleware := web.ContextMiddleware(func(log domain.Log) *web.Context {
		return &web.Context{
			Database: db,
			Log:      log,
		}
	}, context.Log, &logsPersister)

	baseRouter := handler.Base{
		WebContextGetter: web.ContextMiddlewareGetter,
	}

	router := handler.Router{
		User: handler.User{
			Base: baseRouter,
		},
	}

	errorsPersister := repository.Errors{
		Context:          context,
		FallbackFilePath: "./errors",
		Database:         db,
	}

	server := web.Gin(context, func(engine *gin.Engine) {

		engine.Use(web.RequestIDMiddleware(func() string {
			return uuid.NewString()
		}))

		engine.Use(web.ErrorPersisterMiddleware(errorsPersister, context.Log, web.RequestIdGetter))

		engine.Use(web.Recover(func(context *gin.Context, err any) {

			internalError, ok := err.(handler.ServerInternalError)
			if !ok {
				context.JSON(http.StatusInternalServerError, err)
			} else {
				context.JSON(http.StatusInternalServerError, internalError.Body)
			}
		}))

		engine.Use(contextMiddleware)

		apiV1 := engine.Group("/api/v1")
		apiV1.GET("/", router.User.GetById)
	})

	server.Serve()
}
