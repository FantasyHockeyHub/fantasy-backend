package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

var (
	NotFoundPlayerStatistic = errors.New("not found statistic by player id")
)

func NewPlayersService(storage PlayersStorage) *PlayersService {
	return &PlayersService{
		storage: storage,
	}
}

type PlayersStorage interface {
	CreatePlayers(playersData []players.Player) error
	GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error)
	GetPlayerByID(playerID int) (players.PlayerResponse, error)
	GetPlayerCards(filter players.PlayerCardsFilter) ([]players.PlayerCardResponse, error)
	AddPlayerCards(tx *sqlx.Tx, buy store.BuyProductModel) error
	CardUnpacking(id int, userID uuid.UUID) error
	InsertPlayerCards(tx *sqlx.Tx, buy store.BuyProductModel, selectedPlayerIDs []int) error
	GetPlayerStatistics(ctx context.Context, playerID int) ([]players.PlayersStatisticDB, error)
}

type PlayersService struct {
	storage PlayersStorage
}

func (s *PlayersService) CreatePlayers(playersData []players.Player) error {

	err := s.storage.CreatePlayers(playersData)
	if err != nil {
		log.Println("Service. CreatePlayers:", err)
		return err
	}

	return nil
}

func (s *PlayersService) GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error) {

	res, err := s.storage.GetPlayers(playersFilter)
	if err != nil {
		log.Println("Service. GetPlayers:", err)
		return res, err
	}
	playerCardsMap := make(map[int][]players.PlayerCardResponse)

	userCards, err := s.storage.GetPlayerCards(players.PlayerCardsFilter{
		ProfileID:        playersFilter.ProfileID,
		League:           playersFilter.League,
		HasUnpackedParam: true,
		Unpacked:         true,
	})
	if err != nil {
		log.Println("Service. GetPlayerCards:", err)
		return res, err
	}

	for _, card := range userCards {
		playerCardsMap[card.PlayerID] = append(playerCardsMap[card.PlayerID], card)
	}

	for i, player := range res {
		res[i].RarityName = store.PlayerCardsRarityTitles[store.ErrCardRarity]
		if cards, ok := playerCardsMap[player.ID]; ok {
			var hasGoldCard bool
			for _, card := range cards {
				if card.Rarity == store.Gold {
					res[i].CardRarity = store.Gold
					hasGoldCard = true
					break
				}
			}

			if !hasGoldCard && len(cards) > 0 {
				res[i].CardRarity = store.Silver
			}

			res[i].RarityName = store.PlayerCardsRarityTitles[res[i].CardRarity]
		}
	}

	return res, nil
}

func (s *PlayersService) GetPlayerCards(filter players.PlayerCardsFilter) ([]players.PlayerCardResponse, error) {

	res, err := s.storage.GetPlayerCards(filter)
	if err != nil {
		log.Println("Service. GetPlayerCards:", err)
		return res, err
	}

	return res, nil
}

func (s *PlayersService) CardUnpacking(id int, userID uuid.UUID) error {

	err := s.storage.CardUnpacking(id, userID)
	if err != nil {
		log.Println("Service. CardUnpacking:", err)
		return err
	}

	return nil
}

func (s *PlayersService) GetStatisticByPlayerId(ctx context.Context, playerId int) ([]players.PlayersStatisticDB, error) {

	playerStatistic, err := s.storage.GetPlayerStatistics(ctx, playerId)
	if err != nil {
		return playerStatistic, fmt.Errorf("GetPlayerStatistics: %v", err)
	}
	if len(playerStatistic) <= 0 {
		return nil, NotFoundPlayerStatistic
	}
	return playerStatistic, nil
}
