package api

import (
	"encoding/json"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
			log.Println("Error decoding response:", err)
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
				player.Position = players.PlayerPosition["Goalie"]
			case "защитник":
				player.Position = players.PlayerPosition["Defensemen"]
			case "нападающий":
				player.Position = players.PlayerPosition["Forward"]
			}

			allPlayersData = append(allPlayersData, player)
		}

		page++
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
	var nhlPlayers []players.Player
	teams := make([]string, 0, len(tournaments.NHLId))
	for key := range tournaments.NHLId {
		teams = append(teams, key)
	}

	for i := 0; i < len(teams); i++ {
		url := fmt.Sprintf("https://api-web.nhle.com/v1/roster/%s/current", teams[i])
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("CreateNHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, StatusResponse{"error"})
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("CreateNHLPlayers:", err)
			ctx.JSON(http.StatusInternalServerError, StatusResponse{"error"})
			return
		}
		defer res.Body.Close()

		var response players.NHLRosterResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			log.Println("CreateNHLPlayers: Error decoding response body:", err)
			ctx.JSON(http.StatusInternalServerError, StatusResponse{"error"})
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
				player.Position = players.PlayerPosition["Goalie"]
			case "D":
				player.Position = players.PlayerPosition["Defensemen"]
			default:
				player.Position = players.PlayerPosition["Forward"]
			}

			nhlPlayers = append(nhlPlayers, player)
		}
	}

	ctx.JSON(http.StatusOK, StatusResponse{"ок"})
}
