package middleware

import (
	"net/http"

	token "github.com/cyzhang39/go_market/auth"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok := ctx.Request.Header.Get("token")
		if tok == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authentication header"})
			ctx.Abort()
			return
		}
		claim, err := token.ValidateTok(tok)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("email", claim.Email)
		ctx.Set("uid", claim.UID)
		ctx.Next()
	}
}
