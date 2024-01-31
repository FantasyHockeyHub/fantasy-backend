package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api Api) SignUp(ctx *gin.Context) {
	//var req Auth
	//if err := ctx.BindJSON(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
	//	return
	//}

	//err := api.auth.Register(ctx, string(req.PwdHash), string(req.Email))
	//if err == auth.NonExistentDomain {
	//	ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
	//	return
	//}
	//if err == auth.UserAlreadyExist {
	//	ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
	//	return
	//}
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, getInternalServerError())
	//	return
	//}

	ctx.AbortWithStatus(http.StatusOK)
}
