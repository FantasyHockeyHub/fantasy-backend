package api

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	err := api.services.User.SignUp(inp)
	if err != nil {
		log.Println("SignUp:", err)
		switch err {
		case service.UserAlreadyExistsError,
			service.InvalidNicknameError,
			service.NicknameTakenError,
			service.PasswordValidationError,
			service.InvalidVerificationCodeError,
			storage.VerificationCodeError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	tokens, err := api.services.User.SignIn(inp)
	if err != nil {
		log.Println("SignIn:", err)
		switch err {
		case storage.UserDoesNotExistError,
			service.IncorrectPasswordError:
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	err := api.services.User.SendVerificationCode(inp.Email)
	if err != nil {
		log.Println("SendVerificationCode:", err)
		switch err {
		case service.UserAlreadyExistsError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	tokens, err := api.services.User.RefreshTokens(inp.RefreshToken)
	if err != nil {
		log.Println("RefreshTokens:", err)
		switch err {
		case service.InvalidRefreshTokenError,
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	err := api.services.User.Logout(inp.RefreshToken)
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

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}

// CheckUserDataExists godoc
// @Summary Существует ли пользователь с указанными параметрами
// @Schemes
// @Description Существует ли уже пользователь с таким email или nickname. Код 200: пользователь с такими данными уже существует, код 404: пользователь с такими данными не найден.
// @Tags user
// @Accept json
// @Produce json
// @Param email query string false "Email пользователя" Example(test@test.test)
// @Param nickname query string false "Nickname пользователя" Example(Qwerty1)
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /user/exists [get]
func (api Api) checkUserDataExists(ctx *gin.Context) {
	var inp user.UserExistsDataInput

	if err := ctx.ShouldBindQuery(&inp); err != nil || (inp.Email == "" && inp.Nickname == "") {
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
		return
	}

	err := api.services.User.CheckUserDataExists(inp)
	if err != nil {
		log.Println("CheckUserDataExists:", err)
		switch err {
		case service.InvalidNicknameError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		case service.UserDoesNotExistError:
			ctx.JSON(http.StatusNotFound, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"Пользователь с указанными параметрами уже существует"})
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
		log.Println("UserInfo:", err)
		return
	}

	userInfo, err := api.services.User.GetUserInfo(userID)
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		log.Println("ChangePassword:", err)
		return
	}

	var changePasswordData = user.ChangePasswordModel{
		ProfileID:   userID,
		OldPassword: inp.OldPassword,
		NewPassword: inp.NewPassword,
	}

	err = api.services.User.ChangePassword(changePasswordData)
	if err != nil {
		log.Println("ChangePassword:", err)
		switch err {
		case service.PasswordValidationError,
			service.IncorrectPasswordError,
			storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	err := api.services.User.ForgotPassword(inp.Email)
	if err != nil {
		log.Println("ForgotPassword:", err)
		switch err {
		case service.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
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
		ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputBodyError))
		return
	}

	err := api.services.User.ResetPassword(inp)
	if err != nil {
		log.Println("ResetPassword:", err)
		switch err {
		case service.PasswordValidationError,
			storage.ResetHashError,
			storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}

// DeleteProfile godoc
// @Summary Удаление профиля
// @Security ApiKeyAuth
// @Schemes
// @Description Удаление профиля пользователя
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 400,401 {object} Error
// @Failure 500 {object} Error
// @Router /user/delete [delete]
func (api Api) deleteProfile(ctx *gin.Context) {
	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		log.Println("DeleteProfile:", err)
		return
	}

	err = api.services.User.DeleteProfile(userID)
	if err != nil {
		log.Println("DeleteProfile:", err)
		switch err {
		case storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}

// getCoinTransactions godoc
// @Summary Получение истории транзакций пользователя
// @Security ApiKeyAuth
// @Schemes
// @Description Получение истории транзакций пользователя по access токену
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} user.CoinTransactionsModel
// @Failure 400,401 {object} Error
// @Failure 500 {object} Error
// @Router /user/transactions [get]
func (api Api) getCoinTransactions(ctx *gin.Context) {
	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		log.Println("GetCoinTransactions:", err)
		return
	}

	transactions, err := api.services.User.GetCoinTransactions(userID)
	if err != nil {
		log.Println("GetCoinTransactions:", err)
		switch err {
		case storage.UserDoesNotExistError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, transactions)
}
