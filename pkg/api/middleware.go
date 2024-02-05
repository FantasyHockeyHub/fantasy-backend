package api

import (
	"errors"
	user_service "github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

func (api Api) userIdentity(ctx *gin.Context) {
	id, err := api.parseAuthHeader(ctx)
	if err != nil {
		log.Println("Authorization:", err)
		ctx.JSON(http.StatusUnauthorized, getUnauthorizedError(err))
		return
	}

	ctx.Set("userID", id)
}

func (api Api) parseAuthHeader(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return "", user_service.AuthHeaderError
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", user_service.InvalidAuthHeaderError
	}

	if len(headerParts[1]) == 0 {
		return "", user_service.EmptyTokenError
	}

	return api.user.Jwt.ParseJWT(headerParts[1])
}

func parseUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, ok := c.Get("userID")
	if !ok {
		return uuid.Nil, errors.New("userID not found in context")
	}

	parsedUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		return uuid.Nil, errors.New("failed to parse userID as UUID")
	}

	return parsedUserID, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
