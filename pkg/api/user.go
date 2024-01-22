package api

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SignUp godoc
// @Summary Регистрация
// @Schemes
// @Description Регистрация нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param data body user.SignUpInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Router /auth/sign-up [post]
func (api Api) SignUp(ctx *gin.Context) {
	var inp user.SignUpInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.SignUp(ctx, inp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

// SignIn godoc
// @Summary Авторизация
// @Schemes
// @Description Авторизация пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param data body user.SignInInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Router /auth/sign-in [post]
func (api Api) SignIn(ctx *gin.Context) {
	var inp user.SignInInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.SignIn(ctx, inp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email,max=64"`
}

// CheckEmailExists godoc
// @Summary Использован ли данный email в сервисе
// @Schemes
// @Description Существует ли уже пользователь с таким email
// @Tags auth
// @Accept json
// @Produce json
// @Param data body EmailInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Router /auth/check-email [post]
func (api Api) CheckEmailExists(ctx *gin.Context) {
	var inp EmailInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.CheckEmailExists(inp.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

type NicknameInput struct {
	Nickname string `json:"nickname" binding:"required,min=4,max=64"`
}

// CheckNicknameExists godoc
// @Summary Использован ли данный nickname в сервисе
// @Schemes
// @Description Существует ли уже пользователь с таким nickname
// @Tags auth
// @Accept json
// @Produce json
// @Param data body NicknameInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Router /auth/check-nickname [post]
func (api Api) CheckNicknameExists(ctx *gin.Context) {
	var inp NicknameInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.CheckNicknameExists(inp.Nickname)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}
