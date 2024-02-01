package api

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	user_service "github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"log"
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
// @Failure 500 {object} Error
// @Router /auth/sign-up [post]
func (api Api) SignUp(ctx *gin.Context) {
	var inp user.SignUpInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.SignUp(inp)
	if err != nil {
		log.Println("SignUp:", err)
		switch err {
		case user_service.UserAlreadyExistsError,
			user_service.InvalidNicknameError,
			user_service.NicknameTakenError,
			user_service.PasswordValidationError,
			user_service.InvalidVerificationCodeError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
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
// @Success 200 {object} user.Tokens
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/sign-in [post]
func (api Api) SignIn(ctx *gin.Context) {
	var inp user.SignInInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	tokens, err := api.user.SignIn(inp)
	if err != nil {
		log.Println("SignIn:", err)
		switch err {
		case storage.UserDoesNotExistError,
			user_service.IncorrectPasswordError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, tokens)
}

type RefreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokens godoc
// @Summary Обновление токенов
// @Schemes
// @Description Обновление access и refresh токенов
// @Tags auth
// @Accept json
// @Produce json
// @Param data body RefreshInput true "Входные параметры"
// @Success 200 {object} user.Tokens
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/refresh-tokens [post]
func (api Api) RefreshTokens(ctx *gin.Context) {
	var inp RefreshInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	tokens, err := api.user.RefreshTokens(inp.RefreshToken)
	if err != nil {
		log.Println("RefreshTokens:", err)
		switch err {
		case user_service.InvalidRefreshTokenError,
			storage.RefreshTokenNotFoundError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, tokens)
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email,max=64"`
}

// CheckEmailExists godoc
// @Summary Использован ли данный email в сервисе
// @Schemes
// @Description Существует ли уже пользователь с таким email
// @Tags user
// @Accept json
// @Produce json
// @Param data body EmailInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/check-email [post]
func (api Api) CheckEmailExists(ctx *gin.Context) {
	var inp EmailInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.CheckEmailExists(inp.Email)
	if err != nil {
		log.Println("CheckEmailExists:", err)
		switch err {
		case user_service.UserAlreadyExistsError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
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
// @Tags user
// @Accept json
// @Produce json
// @Param data body NicknameInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/check-nickname [post]
func (api Api) CheckNicknameExists(ctx *gin.Context) {
	var inp NicknameInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.CheckNicknameExists(inp.Nickname)
	if err != nil {
		log.Println("CheckNicknameExists:", err)
		switch err {
		case user_service.InvalidNicknameError,
			user_service.NicknameTakenError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

// SendVerificationCode godoc
// @Summary Отправка кода подтверждения
// @Schemes
// @Description Отправка письма с кодом для подтверждения email пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param data body EmailInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/email/send-code [post]
func (api Api) SendVerificationCode(ctx *gin.Context) {
	var inp EmailInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.SendVerificationCode(inp.Email)
	if err != nil {
		log.Println("SendVerificationCode:", err)
		switch err {
		case user_service.UserAlreadyExistsError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}
