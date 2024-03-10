package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	mock_service "github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/mocks"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.SignUpInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.SignUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPassword1", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPassword1",
				Code:     111111,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"nickname": "t", "email": "test", "password": "TestPassword1", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "t",
				Email:    "test",
				Password: "TestPassword1",
				Code:     111111,
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.SignUpInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "User already exists",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPassword1", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPassword1",
				Code:     111111,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(service.UserAlreadyExistsError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.UserAlreadyExistsError),
		},
		{
			name:      "Invalid nickname",
			inputBody: `{"nickname": "test1#@$", "email": "test@test.test", "password": "TestPassword1", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "test1#@$",
				Email:    "test@test.test",
				Password: "TestPassword1",
				Code:     111111,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(service.InvalidNicknameError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.InvalidNicknameError),
		},
		{
			name:      "Nickname taken",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPassword1", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPassword1",
				Code:     111111,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(service.NicknameTakenError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.NicknameTakenError),
		},
		{
			name:      "Password is not valid",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPass", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPass",
				Code:     111111,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(service.PasswordValidationError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.PasswordValidationError),
		},
		{
			name:      "Invalid verification code",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPass", "code": 1}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPass",
				Code:     1,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(service.InvalidVerificationCodeError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.InvalidVerificationCodeError),
		},
		{
			name:      "Verification code not found or expired",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPass", "code": 100000}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPass",
				Code:     100000,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(storage.VerificationCodeError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.VerificationCodeError),
		},
		{
			name:      "Service error",
			inputBody: `{"nickname": "test1", "email": "test@test.test", "password": "TestPassword1", "code": 111111}`,
			inputData: user.SignUpInput{
				Nickname: "test1",
				Email:    "test@test.test",
				Password: "TestPassword1",
				Code:     111111,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignUpInput) {
				s.EXPECT().SignUp(inp).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/auth/sign-up", handler.signUp)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.SignInInput, tokens user.Tokens)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.SignInInput
		tokens               user.Tokens
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "test@test.test", "password": "TestPassword1"}`,
			inputData: user.SignInInput{
				Email:    "test@test.test",
				Password: "TestPassword1",
			},
			tokens: user.Tokens{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				RefreshToken: "6f1ccafb88521208c3e32a603733e2a53ac10da7540b1e756dfbc2345a2950b2",
				ExpiresIn:    1708367533,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignInInput, tokens user.Tokens) {
				s.EXPECT().SignIn(inp).Return(tokens, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c","refreshToken":"6f1ccafb88521208c3e32a603733e2a53ac10da7540b1e756dfbc2345a2950b2","expiresIn":1708367533}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"email": "test", "password": "TestPassword1"}`,
			inputData: user.SignInInput{
				Email:    "test",
				Password: "TestPassword1",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.SignInInput, tokens user.Tokens) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "User does not exist",
			inputBody: `{"email": "test@test.test", "password": "TestPassword1"}`,
			inputData: user.SignInInput{
				Email:    "test@test.test",
				Password: "TestPassword1",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignInInput, tokens user.Tokens) {
				s.EXPECT().SignIn(inp).Return(user.Tokens{}, storage.UserDoesNotExistError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.UserDoesNotExistError),
		},
		{
			name:      "Incorrect password",
			inputBody: `{"email": "test@test.test", "password": "TestPassword1"}`,
			inputData: user.SignInInput{
				Email:    "test@test.test",
				Password: "TestPassword1",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignInInput, tokens user.Tokens) {
				s.EXPECT().SignIn(inp).Return(user.Tokens{}, service.IncorrectPasswordError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.IncorrectPasswordError),
		},
		{
			name:      "Service error",
			inputBody: `{"email": "test@test.test", "password": "TestPassword1"}`,
			inputData: user.SignInInput{
				Email:    "test@test.test",
				Password: "TestPassword1",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.SignInInput, tokens user.Tokens) {
				s.EXPECT().SignIn(inp).Return(user.Tokens{}, errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData, testCase.tokens)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/auth/sign-in", handler.signIn)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/sign-in",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_sendVerificationCode(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.EmailInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.EmailInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "test@test.test"}`,
			inputData: user.EmailInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.EmailInput) {
				s.EXPECT().SendVerificationCode(inp.Email).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"email": "test"}`,
			inputData: user.EmailInput{
				Email: "test",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.EmailInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "User already exists",
			inputBody: `{"email": "test@test.test"}`,
			inputData: user.EmailInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.EmailInput) {
				s.EXPECT().SendVerificationCode(inp.Email).Return(service.UserAlreadyExistsError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.UserAlreadyExistsError),
		},
		{
			name:      "Service error",
			inputBody: `{"email": "test@test.test"}`,
			inputData: user.EmailInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.EmailInput) {
				s.EXPECT().SendVerificationCode(inp.Email).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/email/send-code", handler.sendVerificationCode)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/email/send-code",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_refreshTokens(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.RefreshInput, tokens user.Tokens)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.RefreshInput
		tokens               user.Tokens
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			tokens: user.Tokens{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				RefreshToken: "6f1ccafb88521208c3e32a603733e2a53ac10da7540b1e756dfbc2345a2950b2",
				ExpiresIn:    1708367533,
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput, tokens user.Tokens) {
				s.EXPECT().RefreshTokens(inp.RefreshToken).Return(tokens, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c","refreshToken":"6f1ccafb88521208c3e32a603733e2a53ac10da7540b1e756dfbc2345a2950b2","expiresIn":1708367533}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"refreshToken": "testRefreshToken"}`,
			inputData: user.RefreshInput{
				RefreshToken: "testRefreshToken",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.RefreshInput, tokens user.Tokens) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "Invalid refresh token",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput, tokens user.Tokens) {
				s.EXPECT().RefreshTokens(inp.RefreshToken).Return(user.Tokens{}, service.InvalidRefreshTokenError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.InvalidRefreshTokenError),
		},
		{
			name:      "Refresh token not found",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput, tokens user.Tokens) {
				s.EXPECT().RefreshTokens(inp.RefreshToken).Return(user.Tokens{}, storage.RefreshTokenNotFoundError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.RefreshTokenNotFoundError),
		},
		{
			name:      "Service error",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput, tokens user.Tokens) {
				s.EXPECT().RefreshTokens(inp.RefreshToken).Return(user.Tokens{}, errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData, testCase.tokens)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/auth/refresh-tokens", handler.refreshTokens)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/refresh-tokens",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_logout(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.RefreshInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.RefreshInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput) {
				s.EXPECT().Logout(inp.RefreshToken).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"refreshToken": "testRefreshToken"}`,
			inputData: user.RefreshInput{
				RefreshToken: "testRefreshToken",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.RefreshInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "Refresh token not found",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput) {
				s.EXPECT().Logout(inp.RefreshToken).Return(storage.RefreshTokenNotFoundError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.RefreshTokenNotFoundError),
		},
		{
			name:      "Service error",
			inputBody: `{"refreshToken": "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6"}`,
			inputData: user.RefreshInput{
				RefreshToken: "ce95f9a9d1f536c371bc7d94c72e537c21877323f303cf04b7decd8ettt082d6",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.RefreshInput) {
				s.EXPECT().Logout(inp.RefreshToken).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/auth/logout", handler.logout)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/logout",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_checkUserDataExists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.UserExistsDataInput)

	testTable := []struct {
		name                 string
		inputData            user.UserExistsDataInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK. Email is already taken",
			inputData: user.UserExistsDataInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Пользователь с указанными параметрами уже существует"}`,
		},
		{
			name: "OK. Nickname is already taken",
			inputData: user.UserExistsDataInput{
				Nickname: "testNickname1",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Пользователь с указанными параметрами уже существует"}`,
		},
		{
			name: "Email. Wrong input",
			inputData: user.UserExistsDataInput{
				Email: "test",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.UserExistsDataInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputParametersError),
		},
		{
			name: "Nickname. Wrong input",
			inputData: user.UserExistsDataInput{
				Nickname: "t",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.UserExistsDataInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputParametersError),
		},
		{
			name: "Invalid nickname",
			inputData: user.UserExistsDataInput{
				Nickname: "testNickname1#$",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(service.InvalidNicknameError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.InvalidNicknameError),
		},
		{
			name: "404. Email is not taken",
			inputData: user.UserExistsDataInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(service.UserDoesNotExistError)
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.UserDoesNotExistError),
		},
		{
			name: "404. Nickname is not taken",
			inputData: user.UserExistsDataInput{
				Nickname: "testNickname1",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(service.UserDoesNotExistError)
			},
			expectedStatusCode: 404,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.UserDoesNotExistError),
		},
		{
			name: "Email. Service error",
			inputData: user.UserExistsDataInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
		{
			name: "Nickname. Service error",
			inputData: user.UserExistsDataInput{
				Nickname: "testNickname1",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.UserExistsDataInput) {
				s.EXPECT().CheckUserDataExists(inp).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.GET("/user/exists", handler.checkUserDataExists)

			w := httptest.NewRecorder()

			param, paramValue := "", ""
			if testCase.inputData.Email != "" {
				param = "email"
				paramValue = testCase.inputData.Email
			} else if testCase.inputData.Nickname != "" {
				param = "nickname"
				paramValue = testCase.inputData.Nickname
			}
			paramString := param + "=" + paramValue

			req := httptest.NewRequest("GET", "/user/exists?"+paramString, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_userInfo(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp uuid.UUID, userInfoResponse user.UserInfoModel)
	userID, _ := uuid.Parse("6bc57ea9-c881-47d3-a293-b925ff1ddf72")
	registrationDate, _ := time.Parse("2006-01-02T15:04:05.999999Z", "2024-02-07T15:33:13.414997Z")

	testTable := []struct {
		name                 string
		inputData            uuid.UUID
		userInfoResponse     user.UserInfoModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputData: userID,
			userInfoResponse: user.UserInfoModel{
				ProfileID:        userID,
				Nickname:         "Test1",
				DateRegistration: registrationDate,
				PhotoLink:        "https://goo.su/6ksU1Nz",
				Coins:            1000,
				Email:            "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp uuid.UUID, userInfoResponse user.UserInfoModel) {
				s.EXPECT().GetUserInfo(inp).Return(userInfoResponse, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"profileID":"6bc57ea9-c881-47d3-a293-b925ff1ddf72","nickname":"Test1","dateRegistration":"2024-02-07T15:33:13.414997Z","photoLink":"https://goo.su/6ksU1Nz","coins":1000,"email":"test@test.test"}`,
		},
		{
			name:      "User does not exist",
			inputData: userID,
			mockBehavior: func(s *mock_service.MockUser, inp uuid.UUID, userInfoResponse user.UserInfoModel) {
				s.EXPECT().GetUserInfo(inp).Return(user.UserInfoModel{}, storage.UserDoesNotExistError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.UserDoesNotExistError),
		},
		{
			name:      "Service error",
			inputData: userID,
			mockBehavior: func(s *mock_service.MockUser, inp uuid.UUID, userInfoResponse user.UserInfoModel) {
				s.EXPECT().GetUserInfo(inp).Return(user.UserInfoModel{}, errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData, testCase.userInfoResponse)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.GET("/user/info", func(ctx *gin.Context) {
				ctx.Set("userID", userID.String())
			}, handler.userInfo)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/user/info", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_changePassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput)
	userID, _ := uuid.Parse("6bc57ea9-c881-47d3-a293-b925ff1ddf72")

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.ChangePasswordInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"oldPassword": "TestPassword1", "newPassword": "TestPassword2"}`,
			inputData: user.ChangePasswordInput{
				OldPassword: "TestPassword1",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput) {
				s.EXPECT().ChangePassword(user.ChangePasswordModel{
					OldPassword: inpData.OldPassword,
					NewPassword: inpData.NewPassword,
					ProfileID:   userID,
				}).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"oldPassword": "TestPassword1", "newPassword": "Test"}`,
			inputData: user.ChangePasswordInput{
				OldPassword: "TestPassword1",
				NewPassword: "Test",
			},
			mockBehavior:       func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "Password validation error",
			inputBody: `{"oldPassword": "TestPassword1", "newPassword": "testpassword"}`,
			inputData: user.ChangePasswordInput{
				OldPassword: "TestPassword1",
				NewPassword: "testpassword",
			},
			mockBehavior: func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput) {
				s.EXPECT().ChangePassword(user.ChangePasswordModel{
					OldPassword: inpData.OldPassword,
					NewPassword: inpData.NewPassword,
					ProfileID:   userID,
				}).Return(service.PasswordValidationError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.PasswordValidationError),
		},
		{
			name:      "Incorrect password",
			inputBody: `{"oldPassword": "TestPassword1", "newPassword": "TestPassword2"}`,
			inputData: user.ChangePasswordInput{
				OldPassword: "TestPassword1",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput) {
				s.EXPECT().ChangePassword(user.ChangePasswordModel{
					OldPassword: inpData.OldPassword,
					NewPassword: inpData.NewPassword,
					ProfileID:   userID,
				}).Return(service.IncorrectPasswordError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.IncorrectPasswordError),
		},
		{
			name:      "User does not exist",
			inputBody: `{"oldPassword": "TestPassword1", "newPassword": "TestPassword2"}`,
			inputData: user.ChangePasswordInput{
				OldPassword: "TestPassword1",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput) {
				s.EXPECT().ChangePassword(user.ChangePasswordModel{
					OldPassword: inpData.OldPassword,
					NewPassword: inpData.NewPassword,
					ProfileID:   userID,
				}).Return(storage.UserDoesNotExistError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.UserDoesNotExistError),
		},
		{
			name:      "Service error",
			inputBody: `{"oldPassword": "TestPassword1", "newPassword": "TestPassword2"}`,
			inputData: user.ChangePasswordInput{
				OldPassword: "TestPassword1",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inpModel user.ChangePasswordModel, inpData user.ChangePasswordInput) {
				s.EXPECT().ChangePassword(user.ChangePasswordModel{
					OldPassword: inpData.OldPassword,
					NewPassword: inpData.NewPassword,
					ProfileID:   userID,
				}).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mock_service.NewMockUser(c)
			testCase.mockBehavior(userService, user.ChangePasswordModel{
				OldPassword: testCase.inputData.OldPassword,
				NewPassword: testCase.inputData.NewPassword,
				ProfileID:   userID,
			}, testCase.inputData)

			services := &service.Services{User: userService}
			handler := Api{services: services}

			r := gin.New()
			r.PATCH("/user/password/change", func(ctx *gin.Context) {
				ctx.Set("userID", userID.String())
			}, handler.changePassword)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("PATCH", "/user/password/change",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_forgotPassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.EmailInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.EmailInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "test@test.test"}`,
			inputData: user.EmailInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.EmailInput) {
				s.EXPECT().ForgotPassword(inp.Email).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"email": "test"}`,
			inputData: user.EmailInput{
				Email: "test",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.EmailInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "User does not exist",
			inputBody: `{"email": "test@test.test"}`,
			inputData: user.EmailInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.EmailInput) {
				s.EXPECT().ForgotPassword(inp.Email).Return(service.UserDoesNotExistError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.UserDoesNotExistError),
		},
		{
			name:      "Service error",
			inputBody: `{"email": "test@test.test"}`,
			inputData: user.EmailInput{
				Email: "test@test.test",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.EmailInput) {
				s.EXPECT().ForgotPassword(inp.Email).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/user/password/forgot", handler.forgotPassword)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/user/password/forgot",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_resetPassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, inp user.ResetPasswordInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            user.ResetPasswordInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"hash": "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF", "newPassword": "TestPassword2"}`,
			inputData: user.ResetPasswordInput{
				Hash:        "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.ResetPasswordInput) {
				s.EXPECT().ResetPassword(inp).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"hash": "1", "newPassword": "TestPassword2"}`,
			inputData: user.ResetPasswordInput{
				Hash:        "1",
				NewPassword: "TestPassword2",
			},
			mockBehavior:       func(s *mock_service.MockUser, inp user.ResetPasswordInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, InvalidInputBodyError),
		},
		{
			name:      "Password validation error",
			inputBody: `{"hash": "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF", "newPassword": "testpassword"}`,
			inputData: user.ResetPasswordInput{
				Hash:        "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF",
				NewPassword: "testpassword",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.ResetPasswordInput) {
				s.EXPECT().ResetPassword(inp).Return(service.PasswordValidationError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, service.PasswordValidationError),
		},
		{
			name:      "Password reset hash not found or expired",
			inputBody: `{"hash": "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF", "newPassword": "TestPassword2"}`,
			inputData: user.ResetPasswordInput{
				Hash:        "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.ResetPasswordInput) {
				s.EXPECT().ResetPassword(inp).Return(storage.ResetHashError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.ResetHashError),
		},
		{
			name:      "User does not exist",
			inputBody: `{"hash": "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF", "newPassword": "TestPassword2"}`,
			inputData: user.ResetPasswordInput{
				Hash:        "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.ResetPasswordInput) {
				s.EXPECT().ResetPassword(inp).Return(storage.UserDoesNotExistError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.UserDoesNotExistError),
		},
		{
			name:      "Service error",
			inputBody: `{"hash": "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF", "newPassword": "TestPassword2"}`,
			inputData: user.ResetPasswordInput{
				Hash:        "oeMAzzfzT-zDP5fM3kQCQBbQKA0AdXMF",
				NewPassword: "TestPassword2",
			},
			mockBehavior: func(s *mock_service.MockUser, inp user.ResetPasswordInput) {
				s.EXPECT().ResetPassword(inp).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				InternalServerErrorTitle, InternalServerErrorMessage),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputData)

			services := &service.Services{User: user}
			handler := Api{services: services}

			r := gin.New()
			r.PATCH("/user/password/reset", handler.resetPassword)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("PATCH", "/user/password/reset",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
