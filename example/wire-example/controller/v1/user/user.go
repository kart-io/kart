package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProvideUserAPISet = wire.NewSet(NewUserController)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) Users(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "user")
}
