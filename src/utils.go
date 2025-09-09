package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	gen "github.com/cyzhang39/go_market/auth"
	"github.com/cyzhang39/go_market/db"
	"github.com/cyzhang39/go_market/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	products *mongo.Collection = db.CollectionDB(db.Client, "products")
	users    *mongo.Collection = db.CollectionDB(db.Client, "users")
	validate                   = validator.New()
)

func HashPassword(password string) string {
	bts, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Panic(err)
	}
	return string(bts)
}

func Verify(entered string, stored string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(entered))
	isValid := true
	msg := ""
	if err != nil {
		msg = "Incorrect username or password"
		isValid = false
	}
	return isValid, msg
}

func View() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lst []models.Product
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cs, err := products.Find(c, bson.D{{}})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Oops, Something went wrong")
			return
		}
		err = cs.All(c, &lst)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cs.Close(c)
		err = cs.Err()
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "Invalid")
			return
		}

		defer cancel()
		ctx.IndentedJSON(200, lst)
	}
}

func AdminAdd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var prods models.Product
		defer cancel()
		err := ctx.BindJSON(&prods)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		prods.ID = primitive.NewObjectID()
		_, err = products.InsertOne(c, prods)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item"})
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, "Item added successfully.")
	}
}

func Search() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var prod []models.Product
		query := ctx.Query("name")
		if query == "" {
			log.Println("Empty query")
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid empty query"})
			ctx.Abort()
			return
		}

		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := products.Find(c, bson.M{"name": bson.M{"$regex": query}})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Failed to index with given query")
			return
		}

		err = result.All(c, &prod)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "Invalid")
			return
		}
		defer result.Close(c)
		err = result.Err()
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "Invalid")
			return
		}

		defer cancel()
		ctx.IndentedJSON(200, prod)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		err := c.BindJSON(&user)
		// fmt.Println(1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res := validate.Struct(user)
		// fmt.Println(2)
		if res != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": res})
			return
		}

		cnt, err := users.CountDocuments(ctx, bson.M{"email": user.Email})
		// fmt.Println(3)
		if err != nil {
			log.Panic((err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// fmt.Println(4)
		if cnt > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "duplicate user"})
			return
		}

		cnt, err = users.CountDocuments(ctx, bson.M{"phone": user.Phone})
		// fmt.Println(5)
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// fmt.Println(6)
		if cnt > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone number already used"})
			return
		}

		hPassword := HashPassword(*user.Password)
		user.Password = &hPassword
		user.CreateTime, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdateTime, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UID = user.ID.Hex()
		tok, rf, _ := gen.Generate(*user.Email, *user.FirstName, *user.LastName, *user.Phone)
		user.Token = &tok
		user.Refresh = &rf
		user.Cart = make([]models.UserProd, 0)
		user.AddressInfo = make([]models.Address, 0)
		user.Status = make([]models.Order, 0)

		_, err = users.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user profile"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully signed up")

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var found models.User
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err = users.FindOne(ctx, bson.M{"email": user.Email}).Decode(&found)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username or password"})
			return
		}

		isValid, msg := Verify(*user.Password, *found.Password)
		defer cancel()
		if !isValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		// fmt.Println("Verified")
		tok, rf, _ := gen.Generate(*found.Email, *found.FirstName, *found.LastName, found.UID)
		defer cancel()

		gen.UpdateTok(tok, rf, found.UID)
		c.JSON(http.StatusFound, found)

	}
}
