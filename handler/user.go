package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type User struct {
	Base
}

func (u *User) GetById(c *gin.Context) {
	ctx := u.WebContextMust(c)

	ctx.Log.Infow("GetById", "id", 1)

	c.Error(errors.New("a custom error"))
}
