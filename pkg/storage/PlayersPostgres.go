package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"strings"
	"time"
)

var (
	PlayerCardNotFoundError     = errors.New("карточка с указанным id не найдена")
	IncorrectPlayerCardUserID   = errors.New("userID владельца карточки не совпадает с текущим")
	PlayerCardIsAlreadyUnpacked = errors.New("карточка уже была распакована")
)

func (p *PostgresStorage) CreatePlayers(playersData []players.Player) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO players (api_id, position, name, team_id, sweater_number, photo_link, league)
        SELECT $1, $2, $3, t.team_id, $5, $6, $7
        FROM teams t
        WHERE t.api_id = $4 AND t.league = $7`
	for _, player := range playersData {
		_, err := tx.Exec(query, player.ApiID, player.Position, player.Name, player.TeamApiID, player.SweaterNumber, player.Photo, player.League)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *PostgresStorage) GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error) {
	var res []players.PlayerResponse

	query := "SELECT p.id, p.position, p.name, p.team_id, p.sweater_number, p.photo_link, p.league, p.player_cost, t.team_name FROM players p INNER JOIN teams t ON p.team_id = t.team_id WHERE 1=1"

	if len(playersFilter.Players) > 0 {
		query += " AND p.id IN ("
		for i := range playersFilter.Players {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("%d", playersFilter.Players[i])
		}
		query += ")"
	}

	if len(playersFilter.Teams) > 0 {
		query += " AND p.team_id IN ("
		for i := range playersFilter.Teams {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("%d", playersFilter.Teams[i])
		}
		query += ")"
	}

	if playersFilter.Position != 0 {
		query += fmt.Sprintf(" AND p.position = %d", playersFilter.Position)
	}

	if playersFilter.League != 0 {
		query += fmt.Sprintf(" AND p.league = %d", playersFilter.League)
	}

	query = strings.TrimSuffix(query, "AND")

	err := p.db.Select(&res, query)
	if err != nil {
		return res, err
	}

	if res == nil {
		res = []players.PlayerResponse{}
	} else {
		for i := range res {
			res[i].LeagueName = tournaments.LeagueTitles[res[i].League]
			res[i].PositionName = players.PlayerPositionTitles[res[i].Position]
		}
	}

	return res, nil
}

func (p *PostgresStorage) GetPlayerByID(playerID int) (players.PlayerResponse, error) {
	var player players.PlayerResponse

	query := "SELECT p.id, p.position, p.name, p.team_id, p.sweater_number, p.photo_link, p.league, p.player_cost, t.team_name FROM players p INNER JOIN teams t ON p.team_id = t.team_id WHERE p.id = $1"

	err := p.db.Get(&player, query, playerID)
	if err != nil {
		return player, err
	}

	player.LeagueName = tournaments.LeagueTitles[player.League]
	player.PositionName = players.PlayerPositionTitles[player.Position]

	return player, nil
}

func (p *PostgresStorage) GetPlayerCards(filter players.PlayerCardsFilter) ([]players.PlayerCardResponse, error) {
	var res []players.PlayerCardResponse

	query := "SELECT pc.id, pc.profile_id, pc.player_id, pc.rarity, pc.multiply, pc.bonus_metric, pc.unpacked, p.position, p.name, p.team_id, p.sweater_number, p.photo_link, p.league, t.team_name FROM player_cards pc INNER JOIN players p ON pc.player_id = p.id INNER JOIN teams t ON p.team_id = t.team_id WHERE 1=1"

	if filter.League != 0 {
		query += fmt.Sprintf(" AND p.league = %d", filter.League)
	}

	if filter.ProfileID != uuid.Nil {
		query += fmt.Sprintf(" AND pc.profile_id = '%s'", filter.ProfileID.String())
	}

	if filter.Rarity != 0 {
		query += fmt.Sprintf(" AND pc.rarity = %d", filter.Rarity)
	}

	if filter.HasUnpackedParam {
		if filter.Unpacked {
			query += " AND pc.unpacked = true"
		} else if filter.Unpacked == false {
			query += " AND pc.unpacked = false"
		}
	}

	query += " ORDER BY pc.id"

	err := p.db.Select(&res, query)
	if err != nil {
		return res, err
	}

	if res == nil {
		res = []players.PlayerCardResponse{}
	} else {
		for i := range res {
			res[i].RarityName = store.PlayerCardsRarityTitles[res[i].Rarity]
			res[i].LeagueName = tournaments.LeagueTitles[res[i].League]
			res[i].PositionName = players.PlayerPositionTitles[res[i].Position]
			res[i].BonusMetricName = store.BonusMetricTitles[res[i].BonusMetric]
		}
	}

	return res, nil
}

func (p *PostgresStorage) AddPlayerCards(tx *sqlx.Tx, buy store.BuyProductModel) error {
	allPlayers, err := p.GetPlayers(players.PlayersFilter{League: buy.League})
	if err != nil {
		return err
	}

	userCards, err := p.GetPlayerCards(players.PlayerCardsFilter{ProfileID: buy.ProfileID, League: buy.League, Rarity: buy.Rarity})
	if err != nil {
		return err
	}

	allPlayerIDs := getUniquePlayers(allPlayers, userCards)

	if len(allPlayerIDs) == 0 {
		return GetAllCardsError
	}

	selectedPlayerIDs := generateCards(allPlayerIDs)

	cardsCount := buy.PlayerCardsCount
	if len(selectedPlayerIDs) < buy.PlayerCardsCount {
		cardsCount = len(selectedPlayerIDs)
	}

	selectedPlayerIDs = selectedPlayerIDs[:cardsCount]

	err = p.InsertPlayerCards(tx, buy, selectedPlayerIDs)
	if err != nil {
		return err
	}

	return nil
}

func getUniquePlayers(allPlayers []players.PlayerResponse, userCards []players.PlayerCardResponse) map[int]struct{} {
	allPlayerIDs := make(map[int]struct{})
	for _, player := range allPlayers {
		allPlayerIDs[player.ID] = struct{}{}
	}

	userCardPlayerIDs := make(map[int]struct{})
	for _, card := range userCards {
		userCardPlayerIDs[card.PlayerID] = struct{}{}
	}

	for playerID := range userCardPlayerIDs {
		delete(allPlayerIDs, playerID)
	}

	return allPlayerIDs
}

func generateCards(allPlayers map[int]struct{}) []int {
	var selectedPlayerIDs []int

	rand.Seed(time.Now().UnixNano())
	for playerID := range allPlayers {
		selectedPlayerIDs = append(selectedPlayerIDs, playerID)
	}

	rand.Shuffle(len(selectedPlayerIDs), func(i, j int) {
		selectedPlayerIDs[i], selectedPlayerIDs[j] = selectedPlayerIDs[j], selectedPlayerIDs[i]
	})

	return selectedPlayerIDs
}

func (p *PostgresStorage) InsertPlayerCards(tx *sqlx.Tx, buy store.BuyProductModel, selectedPlayerIDs []int) error {
	query := `INSERT INTO player_cards (profile_id, player_id, rarity, multiply, bonus_metric, unpacked) VALUES `
	var valueStrings []string
	var valueArgs []interface{}

	for idx, playerID := range selectedPlayerIDs {
		playerData, err := p.GetPlayerByID(playerID)
		if err != nil {
			return err
		}
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", idx*6+1, idx*6+2, idx*6+3, idx*6+4, idx*6+5, idx*6+6))
		valueArgs = append(valueArgs, buy.ProfileID, playerID, buy.Rarity, store.CardMultiply[buy.Rarity], playerData.Position, false)
	}

	query += strings.Join(valueStrings, ", ")

	_, err := tx.Exec(query, valueArgs...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *PostgresStorage) CardUnpacking(id int, userID uuid.UUID) error {
	card := players.PlayerCardsFilter{}

	err := p.db.QueryRow("SELECT profile_id, unpacked FROM player_cards WHERE id = $1", id).Scan(&card.ProfileID, &card.Unpacked)
	if err == sql.ErrNoRows {
		return PlayerCardNotFoundError
	} else if err != nil {
		return err
	}

	if card.ProfileID != userID {
		return IncorrectPlayerCardUserID
	}

	if card.Unpacked {
		return PlayerCardIsAlreadyUnpacked
	}

	_, err = p.db.Exec("UPDATE player_cards SET unpacked = true WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
