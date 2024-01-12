package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SignUp godoc
// @Summary Регистрация
// @Schemes
// @Description Прямая регистрация нового пользователя в системе
// @Description 403 — пользователь с данным логином уже существует
// @Description 423 — пользователь с недопустимым почтовым доменом
// @Tags auth
// @Accept json
// @Produce json
// @Param data body Auth true "Входные параметры"
// @Success 200
// @Failure 400 {object} Error
// @Router /auth/signup [post]
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
