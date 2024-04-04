package store

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/google/uuid"
)

type CardRarity int8

const (
	ErrCardRarity CardRarity = iota
	Silver
	Gold
)

var PlayerCardsRarity = map[string]CardRarity{
	"Silver": Silver,
	"Gold":   Gold,
}

var PlayerCardsRarityTitles = map[CardRarity]string{
	Silver: "Silver",
	Gold:   "Gold",
}

func (t *CardRarity) GetCardRarityString() string {
	return PlayerCardsRarityTitles[*t]
}

func (t *CardRarity) GetCardRarityId(str string) CardRarity {
	return PlayerCardsRarity[str]
}

type Product struct {
	ID               int                `json:"id" db:"id"`
	ProductName      string             `json:"productName" db:"product_name"`
	Price            int                `json:"price" db:"price"`
	League           tournaments.League `json:"league" db:"league"`
	LeagueName       string             `json:"leagueName"`
	Rarity           CardRarity         `json:"rarity" db:"rarity"`
	RarityName       string             `json:"rarityName"`
	PlayerCardsCount int                `json:"playerCardsCount" db:"player_cards_count"`
	PhotoLink        string             `json:"photoLink" db:"photo_link"`
}

type BuyProductModel struct {
	ID        int `db:"id"`
	ProfileID uuid.UUID
	Coins     int
	Details   string
}
