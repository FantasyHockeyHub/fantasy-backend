package api

import (
	"encoding/json"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// createKHLPlayers godoc
// @Summary Добавление игроков КХЛ
// @Schemes
// @Description Добавление игроков КХЛ в базу данных
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 500 {object} Error
// @Router /players/khl/create [post]
func (api Api) createKHLPlayers(ctx *gin.Context) {
	var allPlayersData []players.Player

	page := 1
	for {
		url := fmt.Sprintf("https://khl.api.webcaster.pro/api/khl_mobile/players_v2.json?page=%d", page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("CreateKHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("CreateKHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}

		decoder := json.NewDecoder(res.Body)
		var playerInfoList []players.KHLPlayerInfo
		err = decoder.Decode(&playerInfoList)
		if err != nil {
			log.Println("CreateKHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}

		if len(playerInfoList) == 0 {
			break
		}

		for _, playerInfo := range playerInfoList {
			player := players.Player{
				ApiID:         playerInfo.Player.ID,
				Name:          playerInfo.Player.Name,
				SweaterNumber: playerInfo.Player.ShirtNumber,
				Photo:         playerInfo.Player.Image,
				TeamApiID:     playerInfo.Player.Team.ID,
				League:        tournaments.Leagues["KHL"],
			}

			switch playerInfo.Player.Role {
			case "вратарь":
				player.Position = players.Goalie
			case "защитник":
				player.Position = players.Defensemen
			case "нападающий":
				player.Position = players.Forward
			}

			allPlayersData = append(allPlayersData, player)
		}

		page++
	}

	err := api.services.Players.CreatePlayers(allPlayersData)
	if err != nil {
		log.Println("CreateKHLPlayers:", err)
		ctx.JSON(http.StatusInternalServerError, getInternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}

// createNHLPlayers godoc
// @Summary Добавление игроков НХЛ
// @Schemes
// @Description Добавление игроков НХЛ в базу данных
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 500 {object} Error
// @Router /players/nhl/create [post]
func (api Api) createNHLPlayers(ctx *gin.Context) {
	var allPlayersData []players.Player
	teams := make([]string, 0, len(tournaments.NHLId))
	for key := range tournaments.NHLId {
		teams = append(teams, key)
	}

	for i := 0; i < len(teams); i++ {
		url := fmt.Sprintf("https://api-web.nhle.com/v1/roster/%s/current", teams[i])
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("CreateNHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("CreateNHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
		defer res.Body.Close()

		var response players.NHLRosterResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			log.Println("CreateNHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}

		for _, playerInfo := range append(append(response.Forwards, response.Defensemen...), response.Goalies...) {
			player := players.Player{
				ApiID:         playerInfo.ID,
				Name:          playerInfo.FirstName.Name + " " + playerInfo.LastName.Name,
				SweaterNumber: playerInfo.Number,
				Photo:         playerInfo.Photo,
				TeamApiID:     tournaments.NHLId[teams[i]],
				League:        tournaments.Leagues["NHL"],
			}

			switch playerInfo.Position {
			case "G":
				player.Position = players.Goalie
			case "D":
				player.Position = players.Defensemen
			default:
				player.Position = players.Forward
			}

			allPlayersData = append(allPlayersData, player)
		}
	}

	err := api.services.Players.CreatePlayers(allPlayersData)
	if err != nil {
		log.Println("CreateNHLPlayers:", err)
		ctx.JSON(http.StatusInternalServerError, getInternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}

// getPlayers godoc
// @Summary Получение списка игроков
// @Schemes
// @Description Получение списка игроков
// @Tags players
// @Accept json
// @Produce json
// @Param profileID query string false "profileID"
// @Param teams query array false "teams"
// @Param position query string false "position" Enums(G, D, F)
// @Param league query string false "league" Enums(NHL, KHL)
// @Success 200 {array} players.PlayerResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /players [get]
func (api Api) getPlayers(ctx *gin.Context) {
	var filterPlayers players.PlayersFilter

	profileFilter := ctx.Query("profileID")
	teamsFilter := ctx.Query("teams")
	leagueFilter := ctx.Query("league")
	positionFilter := ctx.Query("position")

	if profileFilter != "" {
		parsedUserID, err := uuid.Parse(profileFilter)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
		filterPlayers.ProfileID = parsedUserID
	}

	if teamsFilter != "" {
		teamIds := strings.Split(teamsFilter, ",")
		for _, teamId := range teamIds {
			id, err := strconv.Atoi(strings.TrimSpace(teamId))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
				return
			}
			filterPlayers.Teams = append(filterPlayers.Teams, id)
		}
	}

	if leagueFilter != "" {
		switch leagueFilter {
		case "NHL":
			filterPlayers.League = tournaments.NHL
		case "KHL":
			filterPlayers.League = tournaments.KHL
		default:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
	}

	if positionFilter != "" {
		switch positionFilter {
		case "G":
			filterPlayers.Position = players.Goalie
		case "D":
			filterPlayers.Position = players.Defensemen
		case "F":
			filterPlayers.Position = players.Forward
		default:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}

	}

	res, err := api.services.Players.GetPlayers(filterPlayers)
	if err != nil {
		log.Println("GetPlayers:", err)
		ctx.JSON(http.StatusInternalServerError, getInternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// getPlayerPlayerCards godoc
// @Summary Получение списка карточек игроков
// @Schemes
// @Description Получение списка карточек игроков
// @Tags players
// @Accept json
// @Produce json
// @Param profileID query string false "profileID"
// @Param rarity query string false "rarity" Enums(Silver, Gold)
// @Param league query string false "league" Enums(NHL, KHL)
// @Param unpacked query boolean false "unpacked"
// @Success 200 {array} players.PlayerCardResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /players/cards [get]
func (api Api) getPlayerCards(ctx *gin.Context) {
	var filterPlayers players.PlayerCardsFilter
	profileFilter := ctx.Query("profileID")
	rarityFilter := ctx.Query("rarity")
	leagueFilter := ctx.Query("league")
	unpackedFilter := ctx.Query("unpacked")

	if profileFilter != "" {
		parsedUserID, err := uuid.Parse(profileFilter)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
		filterPlayers.ProfileID = parsedUserID
	}

	if rarityFilter != "" {
		switch rarityFilter {
		case "Silver":
			filterPlayers.Rarity = store.Silver
		case "Gold":
			filterPlayers.Rarity = store.Gold
		default:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
	}

	if leagueFilter != "" {
		switch leagueFilter {
		case "NHL":
			filterPlayers.League = tournaments.NHL
		case "KHL":
			filterPlayers.League = tournaments.KHL
		default:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
	}

	if unpackedFilter != "" {
		switch unpackedFilter {
		case "true":
			filterPlayers.Unpacked = true
		case "false":
			filterPlayers.Unpacked = false
		default:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(InvalidInputParametersError))
			return
		}
		filterPlayers.HasUnpackedParam = true
	}

	res, err := api.services.Players.GetPlayerCards(filterPlayers)
	if err != nil {
		log.Println("GetPlayerCards:", err)
		ctx.JSON(http.StatusInternalServerError, getInternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// createNHLPlayers godoc
// @Summary Распаковка карточки игрока
// @Security ApiKeyAuth
// @Schemes
// @Description Распаковка карточки игрока
// @Tags players
// @Accept json
// @Produce json
// @Param id query int true "id карточки" Example(1)
// @Success 200 {object} StatusResponse
// @Failure 400,401 {object} Error
// @Failure 500 {object} Error
// @Router /players/cards/unpack [post]
func (api Api) cardUnpacking(ctx *gin.Context) {
	userID, err := parseUserIDFromContext(ctx)
	if err != nil {
		log.Println("CardUnpacking:", err)
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

	err = api.services.Players.CardUnpacking(id, userID)
	if err != nil {
		log.Println("CardUnpacking:", err)
		switch err {
		case storage.PlayerCardNotFoundError,
			storage.IncorrectPlayerCardUserID,
			storage.PlayerCardIsAlreadyUnpacked:
			ctx.JSON(http.StatusBadRequest, getBadRequestError(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, getInternalServerError())
			return
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}
