package auth

import (
	"context"
	// "fmt"
	"log"
	"os"
	"time"

	"github.com/cyzhang39/go_market/db"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Signature struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.StandardClaims
}

var SECRET = os.Getenv("SECRET_KEY")
var users *mongo.Collection = db.CollectionDB(db.Client, "users")

func Generate(email string, fName string, lName string, uid string) (signed string, refresh string, err error) {
	// fmt.Println(SECRET != "")
	sig := &Signature{
		Email:          email,
		FirstName:      fName,
		LastName:       lName,
		UID:            uid,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()},
	}
	rfSig := &Signature{
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix()},
	}
	// fmt.Println(1)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, sig).SignedString([]byte(SECRET))
	// fmt.Println(token)
	if err != nil {
		return "", "", err
	}

	rfTok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rfSig).SignedString([]byte(SECRET))
	if err != nil {
		log.Panic(err)
		return
	}
	// fmt.Println(3)
	return token, rfTok, err

}

func ValidateTok(signed string) (claims *Signature, msg string) {
	tok, err := jwt.ParseWithClaims(signed, &Signature{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := tok.Claims.(*Signature)
	if !ok {
		msg = "Invalid token"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token has expired"
		return
	}
	return claims, msg

}

func UpdateTok(signed string, refresh string, uid string) {
	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var update primitive.D
	update = append(update, bson.E{Key: "token", Value: signed})
	update = append(update, bson.E{Key: "refresh", Value: refresh})
	upTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update = append(update, bson.E{Key: "updateTime", Value: upTime})
	upsert := true
	idx := bson.M{"uid": uid}
	op := options.UpdateOptions{Upsert: &upsert}
	_, err := users.UpdateOne(c, idx, bson.D{{Key: "$set", Value: update}}, &op)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}

}
