package api

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// getAllProducts godoc
// @Summary Получение товаров из fantasy магазина
// @Schemes
// @Description Получение списка товаров из fantasy магазина
// @Tags store
// @Accept json
// @Produce json
// @Success 200 {array} store.Product
// @Failure 500 {object} Error
// @Router /store/products [get]
func (api Api) getAllProducts(ctx *gin.Context) {

	products, err := api.services.Store.GetAllProducts()
	if err != nil {
		log.Println("GetAllProducts:", err)
		ctx.JSON(http.StatusInternalServerError, getInternalServerError())
		return
	}

	for i := range products {
		products[i].LeagueName = products[i].League.GetLeagueString()
		products[i].RarityName = products[i].Rarity.GetCardRarityString()
	}

	ctx.JSON(http.StatusOK, products)
}

// buyProduct godoc
// @Summary Покупка товара в магазине
// @Security ApiKeyAuth
// @Schemes
// @Description Покупка товара в магазине по id товара
// @Tags store
// @Accept json
// @Produce json
// @Param id query int true "id товара" Example(1)
// @Success 200 {object} StatusResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /store/products/buy [post]
func (api Api) buyProduct(ctx *gin.Context) {
	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		log.Println("BuyProduct:", err)
		return
	}

	var id int
	query := ctx.Request.URL.Query()
	if query.Has("id") {
		id, err = strconv.Atoi(query.Get("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
	}

	var buy = store.BuyProductModel{
		ID:        id,
		ProfileID: userID,
	}

	err = api.services.Store.BuyProduct(buy)
	if err != nil {
		log.Println("BuyProduct:", err)
		switch err {
		case storage.IncorrectProductID,
			storage.NotEnoughCoinsError,
			storage.GetAllCardsError:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}
