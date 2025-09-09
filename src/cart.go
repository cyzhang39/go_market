package src

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cyzhang39/go_market/db"
	"github.com/cyzhang39/go_market/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	products *mongo.Collection
	users    *mongo.Collection
}

func NewApp(products *mongo.Collection, users *mongo.Collection) *App {
	return &App{
		products: products,
		users:    users,
	}
}

func (app *App) CartAdd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pid := ctx.Query("id")
		if pid == "" {
			log.Println("Invalid product id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid product id"))
			return
		}

		uid := ctx.Query("userID")
		if uid == "" {
			log.Println("Invalide user id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
			return
		}

		pHex, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		err = db.CartAdd(c, app.products, app.users, pHex, uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			// return
		}
		ctx.IndentedJSON(200, "Item successfully added")
	}
}

func (app *App) CartRemove() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pid := ctx.Query("id")
		if pid == "" {
			log.Println("Invalid product id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid product id"))
			return
		}

		uid := ctx.Query("userID")
		if uid == "" {
			log.Println("Invalide user id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
			return
		}

		pHex, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		err = db.CartRemove(c, app.products, app.users, pHex, uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "Item removed successfully")
	}
}

func CartGet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("id")

		if uid == "" {
			ctx.Header("Content-Type", "Application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid user ID"})
			ctx.Abort()
			return
		}

		uHex, _ := primitive.ObjectIDFromHex(uid)
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var cart models.User
		err := users.FindOne(c, bson.D{primitive.E{Key: "id", Value: uHex}}).Decode(&cart)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusNotFound, "Invalid user")
			return
		}
		match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "id", Value: uHex}}}}
		ret := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$cart"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$cart.price"}}}}}}
		result, err := users.Aggregate(c, mongo.Pipeline{match, ret, group})
		if err != nil {
			log.Println(err)
		}

		var lst []bson.M
		err = result.All(c, &lst)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range lst {
			ctx.IndentedJSON(200, json["total"])
			ctx.IndentedJSON(200, cart.Cart)
		}

		c.Done()

	}
}

func (app *App) CartBuy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("id")
		if uid == "" {
			log.Panic("Invalid user id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid user id"))
			// return
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := db.CartBuy(c, app.users, uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}
		ctx.IndentedJSON(200, "Order placed successfully")
	}
}

func (app *App) Buy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pid := ctx.Query("id")
		if pid == "" {
			log.Println("Invalid product id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid product id"))
			return
		}

		uid := ctx.Query("userID")
		if uid == "" {
			log.Println("Invalide user id")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
			return
		}

		pHex, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 8 * time.Second)
		defer cancel()

		err = db.Buy(c, app.products, app.users, pHex, uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}
		ctx.IndentedJSON(200, "Order placed successfully")
	}
}
