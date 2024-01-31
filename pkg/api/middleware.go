package api

import (
	user_service "github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (api Api) UserIdentity(ctx *gin.Context) {
	id, err := api.ParseAuthHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, getUnauthorizedError(err))
		return
	}

	ctx.Set("userID", id)
}

func (api Api) ParseAuthHeader(ctx *gin.Context) (string, error) {
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
