package api

import (
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	mock_service "github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/mocks"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_getAllProducts(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStore, productsResponse []store.Product)

	testTable := []struct {
		name                 string
		productsResponse     []store.Product
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			productsResponse: []store.Product{
				{
					ID:               1,
					ProductName:      "Набор серебряных карточек НХЛ",
					Price:            500,
					League:           1,
					LeagueName:       "NHL",
					Rarity:           1,
					RarityName:       "Silver",
					PlayerCardsCount: 5,
					PhotoLink:        "",
				},
			},
			mockBehavior: func(s *mock_service.MockStore, productsResponse []store.Product) {
				s.EXPECT().GetAllProducts().Return(productsResponse, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"productName":"Набор серебряных карточек НХЛ","price":500,"league":1,"leagueName":"NHL","rarity":1,"rarityName":"Silver","playerCardsCount":5,"photoLink":""}]`,
		},
		{
			name: "Service error",
			mockBehavior: func(s *mock_service.MockStore, productsResponse []store.Product) {
				s.EXPECT().GetAllProducts().Return([]store.Product{}, errors.New("something went wrong"))
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

			store := mock_service.NewMockStore(c)
			testCase.mockBehavior(store, testCase.productsResponse)

			services := &service.Services{Store: store}
			handler := Api{services: services}

			r := gin.New()
			r.GET("/store/products", handler.getAllProducts)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/store/products", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_buyProduct(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStore, inp store.BuyProductModel)
	userID, _ := uuid.Parse("6bc57ea9-c881-47d3-a293-b925ff1ddf72")

	testTable := []struct {
		name                 string
		inputData            store.BuyProductModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputData: store.BuyProductModel{
				ID:        1,
				ProfileID: userID,
			},
			mockBehavior: func(s *mock_service.MockStore, inp store.BuyProductModel) {
				s.EXPECT().BuyProduct(inp).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ок"}`,
		},
		{
			name: "Incorrect product id",
			inputData: store.BuyProductModel{
				ID:        1,
				ProfileID: userID,
			},
			mockBehavior: func(s *mock_service.MockStore, inp store.BuyProductModel) {
				s.EXPECT().BuyProduct(inp).Return(storage.IncorrectProductID)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.IncorrectProductID),
		},
		{
			name: "Not enough coins",
			inputData: store.BuyProductModel{
				ID:        1,
				ProfileID: userID,
			},
			mockBehavior: func(s *mock_service.MockStore, inp store.BuyProductModel) {
				s.EXPECT().BuyProduct(inp).Return(storage.NotEnoughCoinsError)
			},
			expectedStatusCode: 400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s","message":"%s"}`,
				BadRequestErrorTitle, storage.NotEnoughCoinsError),
		},
		{
			name: "Service error",
			inputData: store.BuyProductModel{
				ID:        1,
				ProfileID: userID,
			},
			mockBehavior: func(s *mock_service.MockStore, inp store.BuyProductModel) {
				s.EXPECT().BuyProduct(inp).Return(errors.New("something went wrong"))
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

			store := mock_service.NewMockStore(c)
			testCase.mockBehavior(store, testCase.inputData)

			services := &service.Services{Store: store}
			handler := Api{services: services}

			r := gin.New()
			r.POST("/store/products/buy", func(ctx *gin.Context) {
				ctx.Set("userID", userID.String())
			}, handler.buyProduct)

			w := httptest.NewRecorder()

			param, paramValue := "id", strconv.Itoa(testCase.inputData.ID)
			paramString := param + "=" + paramValue

			req := httptest.NewRequest("POST", "/store/products/buy?"+paramString, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
