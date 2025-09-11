package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cyzhang39/go_market/db"
	"github.com/cyzhang39/go_market/models"
)

var (
	products *mongo.Collection = db.CollectionDB(db.Client, "products")
	users    *mongo.Collection = db.CollectionDB(db.Client, "users")
	valrev                     = validator.New()
)

func ReviewRoutes(r *gin.Engine) {
	rt := r.Group("/products")
	rt.POST("/:pid/reviews", MakeReview)
	rt.GET("/:pid/reviews", ListReviews)
}

func CheckPurchase(ctx context.Context, userID, productID primitive.ObjectID) (bool, error) {
	idx := bson.M{"id": userID, "status.cart.id": productID}
	err := users.FindOne(ctx, idx).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func MakeReview(c *gin.Context) {
	userID, ok := GetUID(c)
	if !ok {
		return
	}
	pHex, err := primitive.ObjectIDFromHex(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productId"})
		return
	}

	var body struct {
		Rating float32 `json:"rating" validate:"required,gte=0,lte=5"`
		Review string  `json:"review" validate:"required,min=1,max=4000"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := valrev.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	okBought, err := CheckPurchase(ctx, userID, pHex)
	fmt.Println(okBought)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check purchase failed"})
		return
	}
	if !okBought {
		c.JSON(http.StatusForbidden, gin.H{"error": "only buyers can make reviews"})
		return
	}

	var existing models.Review
	err = db.Reviews.FindOne(ctx, bson.M{"pid": pHex, "uid": userID}).Decode(&existing)
	now := time.Now()

	switch err {
	case nil:
		temp := float64(body.Rating) - float64(existing.Rating)

		idx := bson.M{"id": existing.ID}
		update := bson.M{"$set": bson.M{"rating": body.Rating, "review": body.Review, "updatedAt": now,}}
		_, err = db.Reviews.UpdateOne(ctx, idx, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update review"})
			return
		}

		idx = bson.M{"id": pHex}
		update = bson.M{"$inc": bson.M{"ratingSum": temp}}
		_, err = products.UpdateOne(ctx, idx, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to adjust product rating sum"})
			return
		}

		var prod models.Product
		if err := products.FindOne(ctx, bson.M{"id": pHex}).Decode(&prod); err == nil && prod.RatingCnt > 0 {
			avg := float32(prod.RatingSum / float64(prod.RatingCnt))
			_, _ = products.UpdateOne(ctx, bson.M{"id": pHex}, bson.M{"$set": bson.M{"ratingAvg": avg}})
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})

	case mongo.ErrNoDocuments:
		r := models.Review{
			ID:        primitive.NewObjectID(),
			PID:       pHex,
			UID:       userID,
			Rating:    body.Rating,
			Review:    body.Review,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if _, err := db.Reviews.InsertOne(ctx, r); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
			return
		}

		idx := bson.M{"id": pHex}
		update := bson.M{"$inc": bson.M{"ratingCnt": 1, "ratingSum": float64(body.Rating)}}
		_, err = products.UpdateOne(ctx, idx, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't calculate total review"})
			return
		}

		var prod models.Product
		if err := products.FindOne(ctx, bson.M{"id": pHex}).Decode(&prod); err == nil && prod.RatingCnt > 0 {
			avg := float32(prod.RatingSum / float64(prod.RatingCnt))
			_, _ = products.UpdateOne(ctx, bson.M{"id": pHex}, bson.M{"$set": bson.M{"ratingAvg": avg}})
		}

		c.JSON(http.StatusOK, gin.H{"status": "created"})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "review not found"})
	}
}

func ListReviews(c *gin.Context) {
	pHex, err := primitive.ObjectIDFromHex(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productId"})
		return
	}
	limit := int64(50)
	if v := c.DefaultQuery("limit", "50"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := db.Reviews.Find(ctx, bson.M{"pid": pHex}, options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}}).SetLimit(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't find review"})
		return
	}
	var out []models.Review
	if err := cur.All(ctx, &out); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid cursor read"})
		return
	}
	c.JSON(http.StatusOK, out)
}
