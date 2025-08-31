package routes

import (
	"github.com/cyzhang39/go_market/src"
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {
	route.POST("/users/signup", src.Signup())
	route.POST("/users/login", src.Login())
	route.GET("/users/view", src.View())
	route.GET("/users/search", src.Search())
	route.POST("/admin/product", src.AdminAdd())

}
