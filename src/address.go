package src

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cyzhang39/go_market/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddressAdd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("id")
		if uid == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid uid"})
			ctx.Abort()
			return
		}
		uHex, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Internal server error")

		}
		var address models.Address

		address.ID = primitive.NewObjectID()
		err = ctx.BindJSON(&address)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "id", Value: uHex}}}}
		ret := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$addressInfo"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		result, err := users.Aggregate(c, mongo.Pipeline{match, ret, group})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}
		var info []bson.M
		err = result.All(c, &info)
		if err != nil {
			panic(err)
		}

		var size int32
		for _, num := range info {
			cnt := num["count"]
			size = cnt.(int32)
		}
		if size <= 1 {
			idx := bson.D{primitive.E{Key: "id", Value: uHex}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "addressInfo", Value: address}}}}
			_, err = users.UpdateOne(c, idx, update)
			if err != nil {
				fmt.Println(err)
			}
			ctx.IndentedJSON(200, "Address added")
		} else {
			ctx.IndentedJSON(400, "Action not allowed")
		}
		defer cancel()
		c.Done()

	}
}

func HomeEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("id")
		if uid == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid empty ID"})
			ctx.Abort()
			return
		}

		uHex, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}
		var edt models.Address
		err = ctx.BindJSON(&edt)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			// return
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		idx := bson.D{primitive.E{Key: "id", Value: uHex}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addressInfo.0.house", Value: edt.House}, {Key: "addressInfo.0.street", Value: edt.Street}, {Key: "addressInfo.0.city", Value: edt.City}, {Key: "addressInfo.0.postal", Value: edt.Postal}}}}
		_, err = users.UpdateOne(c, idx, update)
		if err != nil {
			ctx.IndentedJSON(500, "Oops, something went wrong")
			return
		}
		defer cancel()
		c.Done()
		ctx.IndentedJSON(200, "Home address updated")

	}
}

func WorkEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("id")
		if uid == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid empty ID"})
			ctx.Abort()
			return
		}

		uHex, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}
		var edt models.Address
		err = ctx.BindJSON(&edt)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			// return
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		idx := bson.D{primitive.E{Key: "id", Value: uHex}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addressInfo.1.house", Value: edt.House}, {Key: "addressInfo.1.street", Value: edt.Street}, {Key: "addressInfo.1.city", Value: edt.City}, {Key: "addressInfo.1.postal", Value: edt.Postal}}}}
		_, err = users.UpdateOne(c, idx, update)
		if err != nil {
			ctx.IndentedJSON(500, "Oops, something went wrong")
			return
		}
		defer cancel()
		c.Done()
		ctx.IndentedJSON(200, "Work address updated")
	}
}

func AddressDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("id")
		if uid == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid empty ID"})
			ctx.Abort()
			return
		}
		address := make([]models.Address, 0)
		uHex, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		idx := bson.D{primitive.E{Key: "id", Value: uHex}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addressInfo", Value: address}}}}
		_, err = users.UpdateOne(c, idx, update)
		if err != nil {
			ctx.IndentedJSON(400, "Invalid")
			return
		}
		defer cancel()
		c.Done()
		ctx.IndentedJSON(200, "Address deleted")
	}
}
