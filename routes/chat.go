package routes

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cyzhang39/go_market/db"
	"github.com/cyzhang39/go_market/models"
)

var valchat = validator.New()

func ChatRoutes(r *gin.Engine) {
	rt := r.Group("/chats")
	rt.POST("", CreateChat)
	rt.GET("", ListChats)
	rt.POST("/:chatId/messages", SendMsg)
	rt.GET("/:chatId/messages", ListMsg)
	rt.POST("/:chatId/read", Read)
}

func GetUID(c *gin.Context) (primitive.ObjectID, bool) {
	uid := c.Query("userID")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return primitive.NilObjectID, false
	}
	uHex, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return primitive.NilObjectID, false
	}
	return uHex, true
}

func CreateChat(c *gin.Context) {
	curr, ok := GetUID(c)
	if !ok {
		return
	}
	var peer struct {
		PeerID string `json:"peerId" validate:"required,len=24"`
	}
	if err := c.BindJSON(&peer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := valchat.Struct(peer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := primitive.ObjectIDFromHex(peer.PeerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Id"})
		return
	}
	if p == curr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat with yourself not allowed"})
		return
	}

	members := []primitive.ObjectID{curr, p}
	sort.Slice(members, func(i, j int) bool { return members[i].Hex() < members[j].Hex()})

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	var chat models.Chat
	err = db.Chats.FindOne(ctx, bson.M{"members": members}).Decode(&chat)
	if err == nil {
		c.JSON(http.StatusOK, chat)
		return
	}

	now := time.Now()
	chat = models.Chat{
		ID:        primitive.NewObjectID(),
		Members:   members,
		CreatedAt: now,
		UpdatedAt: now,
		UnreadBy:  map[string]int{curr.Hex(): 0, p.Hex(): 0},
	}
	if _, err := db.Chats.InsertOne(ctx, chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	c.JSON(http.StatusCreated, chat)
}

func ListChats(c *gin.Context) {
	curr, ok := GetUID(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	temp, err := db.Chats.Find(ctx, bson.M{"members": curr}, options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var chats []models.Chat
	if err := temp.All(ctx, &chats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chats)
}

func SendMsg(c *gin.Context) {
	curr, ok := GetUID(c)
	if !ok {
		return
	}

	chatId := c.Param("chatId")
	chatHex, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chatId"})
		return
	}

	var body struct {
		Text string `json:"text" validate:"required,min=1,max=4000"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := valchat.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Chats.FindOne(ctx, bson.M{"id": chatHex, "members": curr}).Err(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a participant or chat not found"})
		return
	}

	msg := models.Message{
		ID:        primitive.NewObjectID(),
		ChatID:    chatHex,
		SenderID:  curr,
		Text:      body.Text,
		CreatedAt: time.Now(),
		ReadBy:    []primitive.ObjectID{curr},
	}
	if _, err := db.Messages.InsertOne(ctx, msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})
		return
	}

	var chat models.Chat
	if err := db.Chats.FindOne(ctx, bson.M{"id": chatHex}).Decode(&chat); err == nil {
		other := GetOther(chat.Members, curr)
		idx := bson.M{"id": chatHex}
		update := bson.M{
			"$set": bson.M{"updatedAt": time.Now(), "lastMessage": models.MessagePreview{Text: msg.Text, SenderID:  msg.SenderID, CreatedAt: msg.CreatedAt,}}, 
			"$inc": bson.M{"unreadBy." + other.Hex(): 1},
		}
		_, err = db.Chats.UpdateOne(ctx, idx, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, failed to send message"})
		}
	}

	c.JSON(http.StatusCreated, msg)
}

func ListMsg(c *gin.Context) {
	curr, ok := GetUID(c)
	if !ok {
		return
	}

	chatId := c.Param("chatId")
	chatHex, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chatId"})
		return
	}

	limit := Limit(c.DefaultQuery("limit", "50"), 1, 200)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Chats.FindOne(ctx, bson.M{"id": chatHex, "members": curr}).Err(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a participant or chat not found"})
		return
	}

	cur, err := db.Messages.Find(
		ctx,
		bson.M{"chatId": chatHex},
		options.Find().
			SetSort(bson.D{{Key: "createdAt", Value: -1}}).
			SetLimit(limit),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	var msgs []models.Message
	if err := cur.All(ctx, &msgs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cursor read failed"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}

func Read(c *gin.Context) {
	curr, ok := GetUID(c)
	if !ok {
		return
	}

	chatId := c.Param("chatId")
	chatHex, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chatId"})
		return
	}

	var temp struct {
		Upto string `json:"upto"`
	}
	_ = c.BindJSON(&temp)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.Chats.FindOne(ctx, bson.M{"id": chatHex, "members": curr}).Err(); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a participant or chat not found"})
		return
	}

	filter := bson.M{"chatId": chatHex}
	if temp.Upto != "" {
		if uptoID, err := primitive.ObjectIDFromHex(temp.Upto); err == nil {
			var uptoMsg models.Message
			if err := db.Messages.FindOne(ctx, bson.M{"id": uptoID, "chatId": chatHex}).Decode(&uptoMsg); err == nil {
				filter["createdAt"] = bson.M{"$lte": uptoMsg.CreatedAt}
			}
		}
	}

	if _, err := db.Messages.UpdateMany(ctx, filter, bson.M{"$addToSet": bson.M{"readBy": curr}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark read"})
		return
	}

	idx := bson.M{"id": chatHex}
	update := bson.M{"$set": bson.M{"unreadBy." + curr.Hex(): 0}}
	_, err = db.Chats.UpdateOne(ctx, idx, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, could not mark as read"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetOther(members []primitive.ObjectID, curr primitive.ObjectID) primitive.ObjectID {
	if members[0] == curr {
		return members[1]
	}
	return members[0]
}

func Limit(s string, min, max int64) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
