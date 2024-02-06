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
func (api Api) signUp(ctx *gin.Context) {
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
			user_service.InvalidVerificationCodeError,
			storage.VerificationCodeError:
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
func (api Api) signIn(ctx *gin.Context) {
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

// SendVerificationCode godoc
// @Summary Отправка кода подтверждения
// @Schemes
// @Description Отправка письма с кодом для подтверждения email пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param data body user.EmailInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/email/send-code [post]
func (api Api) sendVerificationCode(ctx *gin.Context) {
	var inp user.EmailInput
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

// RefreshTokens godoc
// @Summary Обновление токенов
// @Schemes
// @Description Обновление access и refresh токенов
// @Tags auth
// @Accept json
// @Produce json
// @Param data body user.RefreshInput true "Входные параметры"
// @Success 200 {object} user.Tokens
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/refresh-tokens [post]
func (api Api) refreshTokens(ctx *gin.Context) {
	var inp user.RefreshInput
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

// Logout godoc
// @Summary Выход из системы
// @Schemes
// @Description Выход пользователя из системы
// @Tags auth
// @Accept json
// @Produce json
// @Param data body user.RefreshInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/logout [post]
func (api Api) logout(ctx *gin.Context) {
	var inp user.RefreshInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.Logout(inp.RefreshToken)
	if err != nil {
		log.Println("Logout:", err)
		switch err {
		case storage.RefreshTokenNotFoundError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

// CheckEmailExists godoc
// @Summary Использован ли данный email в сервисе
// @Schemes
// @Description Существует ли уже пользователь с таким email
// @Tags user
// @Accept json
// @Produce json
// @Param data body user.EmailInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/check-email [post]
func (api Api) checkEmailExists(ctx *gin.Context) {
	var inp user.EmailInput
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

// CheckNicknameExists godoc
// @Summary Использован ли данный nickname в сервисе
// @Schemes
// @Description Существует ли уже пользователь с таким nickname
// @Tags user
// @Accept json
// @Produce json
// @Param data body user.NicknameInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/check-nickname [post]
func (api Api) checkNicknameExists(ctx *gin.Context) {
	var inp user.NicknameInput
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

// UserInfo godoc
// @Summary Получение информации о пользователе
// @Security ApiKeyAuth
// @Schemes
// @Description Получение пользовательской информации по access токену
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} user.UserInfoModel
// @Failure 400,401 {object} Error
// @Failure 500 {object} Error
// @Router /user/info [get]
func (api Api) userInfo(ctx *gin.Context) {
	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		return
	}

	userInfo, err := api.user.GetUserInfo(userID)
	if err != nil {
		log.Println("UserInfo:", err)
		switch err {
		case storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, userInfo)
}

// ChangePassword godoc
// @Summary Смена пароля
// @Security ApiKeyAuth
// @Schemes
// @Description Смена пароля
// @Tags user
// @Accept json
// @Produce json
// @Param data body user.ChangePasswordInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400,401 {object} Error
// @Failure 500 {object} Error
// @Router /user/password/change [patch]
func (api Api) changePassword(ctx *gin.Context) {
	var inp user.ChangePasswordInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		return
	}

	var changePasswordData = user.ChangePasswordModel{
		ProfileID:   userID,
		OldPassword: inp.OldPassword,
		NewPassword: inp.NewPassword,
	}

	err = api.user.ChangePassword(changePasswordData)
	if err != nil {
		log.Println("ChangePassword:", err)
		switch err {
		case user_service.PasswordValidationError,
			user_service.IncorrectPasswordError,
			storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

// ForgotPassword godoc
// @Summary Забыли пароль
// @Schemes
// @Description Отправка письма с информацией о сбросе пароля
// @Tags user
// @Accept json
// @Produce json
// @Param data body user.EmailInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/password/forgot [post]
func (api Api) forgotPassword(ctx *gin.Context) {
	var inp user.EmailInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.ForgotPassword(inp.Email)
	if err != nil {
		log.Println("ForgotPassword:", err)
		switch err {
		case user_service.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

// ResetPassword godoc
// @Summary Восстановление пароля
// @Schemes
// @Description Восстановление пароля
// @Tags user
// @Accept json
// @Produce json
// @Param data body user.ResetPasswordInput true "Входные параметры"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /user/password/reset [patch]
func (api Api) resetPassword(ctx *gin.Context) {
	var inp user.ResetPasswordInput
	if err := ctx.BindJSON(&inp); err != nil {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
		return
	}

	err := api.user.ResetPassword(inp)
	if err != nil {
		log.Println("ResetPassword:", err)
		switch err {
		case user_service.PasswordValidationError,
			storage.ResetHashError,
			storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}
