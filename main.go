package main

import (
	"log"
	"os"

	"github.com/cyzhang39/go_market/db"
	"github.com/cyzhang39/go_market/middleware"
	"github.com/cyzhang39/go_market/routes"
	"github.com/cyzhang39/go_market/src"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := src.NewApp(db.CollectionDB(db.Client, "products"), db.CollectionDB(db.Client, "users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.Routes(router)
	router.Use(middleware.Authenticate())
	router.GET("/add", server.CartAdd())
	router.GET("/remove", server.CartRemove())
	router.GET("/list", src.CartGet())
	router.GET("/checkout", server.CartBuy())
	router.GET("/buy", server.Buy())
	router.POST("/addressadd", src.AddressAdd())
	router.PUT("/addresshomeedit", src.HomeEdit())
	router.PUT("/addressworkedit", src.WorkEdit())
	router.GET("/addressdel", src.AddressDelete())

	log.Fatal(router.Run(":" + port))
}
