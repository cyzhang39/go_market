package db

import (
	"context"
	"log"
	"time"
	"errors"
	"github.com/cyzhang39/go_market/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidProduct = errors.New("invalid Product")
	ErrInvalidUser = errors.New("invalid User")
	ErrInvalidCart = errors.New("unable to process cart action")
)

func CartAdd(ctx context.Context, products *mongo.Collection, users *mongo.Collection, pid primitive.ObjectID, uid string) error {
	search, err := products.Find(ctx, bson.M{"id": pid})
	if err != nil {
		log.Println(err)
		return ErrInvalidCart
	}

	var uProd []models.UserProd
	err = search.All(ctx, &uProd)
	if err != nil {
		log.Println(err)
		return ErrInvalidProduct
	}

	uHex, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		log.Println(err)
		return ErrInvalidUser
	}

	idx := bson.D{primitive.E{Key: "id", Value: uHex}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "cart", Value: bson.D{{Key: "$each", Value: uProd}}}}}}

	_, err = users.UpdateOne(ctx, idx, update)
	if err != nil {
		return ErrInvalidUser
	}
	return nil 

}

func CartRemove(ctx context.Context, products *mongo.Collection, users *mongo.Collection, pid primitive.ObjectID, uid string) error {
	uHex, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		log.Println(err)
		return ErrInvalidUser
	}

	idx := bson.D(primitive.D{primitive.E{Key: "id", Value: uHex}})
	update := bson.M{"$pull": bson.M{"cart": bson.M{"id": pid}}}

	_, err = users.UpdateMany(ctx, idx, update)
	if err != nil {
		return ErrInvalidCart
	}
	return nil
}

// func CartGet() {

// }

func CartBuy(ctx context.Context, users *mongo.Collection, uid string) error {
	uHex, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		log.Println(err)
		return ErrInvalidUser
	}
	var user models.User
	var order models.Order

	order.ID = primitive.NewObjectID()
	order.OrderTime = time.Now()
	order.Cart = make([]models.UserProd, 0)
	order.Payment.Cash = true

	ret := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$cart"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$cart.price"}}}}}}
	res, err := users.Aggregate(ctx, mongo.Pipeline{ret, group})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var uCart []bson.M
	err = res.All(ctx, &uCart)
	if err != nil {
		panic(err)
	}
	var price float64
	for _, item := range uCart {
		p := item["total"]
		price = p.(float64)

	}
	order.Price = float64(price)

	idx := bson.D{primitive.E{Key: "id", Value: uHex}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order}}}}
	_, err = users.UpdateMany(ctx, idx, update)
	if err != nil {
		log.Println(err)

	}
	err = users.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: uHex}}).Decode(&user)
	if err != nil {
		log.Println(err)

	}

	idx2 := bson.D{primitive.E{Key: "id", Value: uHex}}
	update2 := bson.M{"$push": bson.M{"orders.$[].cart": bson.M{"$each": user.Cart}}}
	_, err = users.UpdateOne(ctx, idx2, update2)
	if err != nil {
		log.Println(err)
	}

	empty := make([]models.UserProd, 0)
	idx3 := bson.D{primitive.E{Key: "id", Value: uHex}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "cart", Value: empty}}}}
	_, err = users.UpdateOne(ctx, idx3, update3)
	if err != nil {
		log.Println(err)
	}
	return nil

}

func Buy(ctx context.Context, products *mongo.Collection, users *mongo.Collection, pid primitive.ObjectID, uid string) error {
	uHex, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		log.Println(err)
		return ErrInvalidUser
	}

	var uProd models.UserProd
	var order models.Order
	order.ID = primitive.NewObjectID()
	order.OrderTime = time.Now()
	order.Cart = make([]models.UserProd, 0)
	order.Payment.Cash = true
	
	err = products.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: pid}}).Decode(&uProd)
	if err != nil {
		log.Println(err)
	}
	
	order.Price = uProd.Price
	idx := bson.D{primitive.E{Key: "id", Value: uHex}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order}}}}
	_, err = users.UpdateOne(ctx, idx, update)
	if err != nil {
		log.Println(err)
	}

	idx2 := bson.D{primitive.E{Key: "id", Value: uHex}}
	update2 := bson.M{"$push": bson.M{"orders.$[].cart": uProd}}
	_, err = users.UpdateOne(ctx, idx2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}