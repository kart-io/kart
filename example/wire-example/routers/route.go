package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kart-io/kart/example/wire-example/controller"
)

func InitRoute(g *gin.Engine, ctr *controller.ApiController) {
	g.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "test")
	})

	// users
	g.GET("/users", ctr.UserController.Users)
}
