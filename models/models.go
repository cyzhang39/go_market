package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	FirstName   *string            `json:"firstName" validate:"required,min=1,max=25"`
	LastName    *string            `json:"lastName" validate:"required,min=1,max=25"`
	Password    *string            `json:"password" validate:"required,min=8,max=32"`
	Email       *string            `json:"email" validate:"required"`
	Phone       *string            `json:"phone"`
	Token       *string            `json:"token"`
	Refresh      *string            `json:"refresh"`
	CreateTime  time.Time          `json:"createTime"`
	UpdateTime  time.Time          `json:"updateTime"`
	UID         string             `json:"uid"`
	Cart        []UserProd         `json:"cart" bson:"cart"`
	AddressInfo []Address          `json:"addressInfo" bson:"addressInfo"`
	Status      []Order            `json:"status" bson:"status"`
}

type Product struct {
	ID     primitive.ObjectID `bson:"id"`
	Name   *string            `json:"name"`
	Price  *float32           `json:"price"`
	Rating *float32           `json:"rating"`
	Img    *string            `json:"img"`
}

type UserProd struct {
	ID     primitive.ObjectID `bson:"id"`
	Name   *string            `json:"name" bson:"name"`
	Price  float32            `json:"price" bson:"price"`
	Rating *float32           `json:"rating" bson:"rating"`
	Img    *string            `json:"img" bson:"img"`
}

type Address struct {
	ID     primitive.ObjectID `bson:"id"`
	House  *string            `json:"house" bson:"house"`
	Street *string            `json:"street" bson:"street"`
	City   *string            `json:"city" bson:"city"`
	Postal *uint8             `json:"postal" bson:"postal"`
}

type Order struct {
	ID        primitive.ObjectID `bson:"id"`
	Cart      []UserProd         `json:"cart" bson:"cart"`
	OrderTime time.Time          `json:"orderTime" bson:"orderTime"`
	Price     float32            `json:"price" bson:"price"`
	DC        *float32           `json:"dc" bson:"dc"`
	Payment   Payment            `json:"payment" bson:"payment"`
}

type Payment struct {
	Online bool
	Cash   bool
}
