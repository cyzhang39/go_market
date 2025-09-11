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
	Email       *string            `json:"email" validate:"email,required"`
	Phone       *string            `json:"phone"`
	Verified    bool               `json:"verified" bson:"verified"`
	Code        string             `json:"code" bson:"code"`
	Token       *string            `json:"token"`
	Refresh     *string            `json:"refresh"`
	CreateTime  time.Time          `json:"createTime"`
	UpdateTime  time.Time          `json:"updateTime"`
	UID         string             `json:"uid"`
	Cart        []UserProd         `json:"cart" bson:"cart"`
	AddressInfo []Address          `json:"addressInfo" bson:"addressInfo"`
	Status      []Order            `json:"status" bson:"status"`
}

type Verification struct {
	Email string `json:"email" validate:"email,required"`
	Code  string `json:"code" validate:"required,len=6"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"id"`
	Name        *string            `json:"name"`
	Price       *float64           `json:"price" validate:"gte=0"`
	Img         *string            `json:"img"`
	Description *string            `json:"description" bson:"description"`
	RatingAvg   float32            `json:"ratingAvg" bson:"ratingAvg"`
	RatingCnt   int64              `json:"ratingCnt" bson:"ratingCnt"`
	RatingSum   float64            `json:"ratingSum" bson:"ratingSum"`
}

type UserProd struct {
	ID     primitive.ObjectID `bson:"id"`
	Name   *string            `json:"name" bson:"name"`
	Price  float64            `json:"price" bson:"price"`
	Rating *float32           `json:"rating" bson:"rating"`
	Img    *string            `json:"img" bson:"img"`
}

type Address struct {
	ID     primitive.ObjectID `bson:"id"`
	House  *string            `json:"house" bson:"house"`
	Street *string            `json:"street" bson:"street"`
	City   *string            `json:"city" bson:"city"`
	Postal *string            `json:"postal" bson:"postal"`
}

type Order struct {
	ID        primitive.ObjectID `bson:"id"`
	Cart      []UserProd         `json:"cart" bson:"cart"`
	OrderTime time.Time          `json:"orderTime" bson:"orderTime"`
	Price     float64            `json:"price" bson:"price"`
	DC        *float32           `json:"dc" bson:"dc"`
	Payment   Payment            `json:"payment" bson:"payment"`
}

type Payment struct {
	Online bool
	Cash   bool
}

type Chat struct {
	ID          primitive.ObjectID   `json:"id" bson:"id"`
	Members     []primitive.ObjectID `json:"members" bson:"members"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt" bson:"updatedAt"`
	LastMessage *MessagePreview      `json:"lastMessage" bson:"lastMessage"`
	UnreadBy    map[string]int       `json:"unreadBy" bson:"unreadBy"`
}

type MessagePreview struct {
	Text      string             `json:"text" bson:"text"`
	SenderID  primitive.ObjectID `json:"senderId" bson:"senderId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type Message struct {
	ID        primitive.ObjectID   `json:"id" bson:"id"`
	ChatID    primitive.ObjectID   `json:"chatId" bson:"chatId"`
	SenderID  primitive.ObjectID   `json:"senderId" bson:"senderId"`
	Text      string               `json:"text" bson:"text" validate:"required,min=1,max=4000"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	ReadBy    []primitive.ObjectID `json:"readBy" bson:"readBy"`
}

type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"id"`
	PID       primitive.ObjectID `json:"pid" bson:"pid"`
	UID       primitive.ObjectID `json:"uid" bson:"uid"`
	Rating    float32            `json:"rating" bson:"rating" validate:"gte=0,lte=5"`
	Review    string             `json:"review" bson:"review" validate:"required,min=1,max=4000"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
