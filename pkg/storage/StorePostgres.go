package storage

import (
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
)

var (
	IncorrectProductID = errors.New("некорректный id товара")
)

func (p *PostgresStorage) GetAllProducts() ([]store.Product, error) {
	var products []store.Product

	err := p.db.Select(&products, `SELECT id, product_name, price, league, rarity, player_cards_count, photo_link
														FROM fantasy_store`)
	if err != nil {
		return products, err
	}

	if products == nil {
		products = []store.Product{}
	} else {
		for i := range products {
			products[i].LeagueName = products[i].League.GetLeagueString()
			products[i].RarityName = products[i].Rarity.GetCardRarityString()
		}
	}

	return products, nil
}

func (p *PostgresStorage) GetProductByID(id int) (store.Product, error) {
	var product store.Product

	err := p.db.Get(&product, `SELECT id, product_name, price, league, rarity, player_cards_count, photo_link
														FROM fantasy_store WHERE id = $1`, id)
	if err != nil {
		return product, err
	}

	if product.ID == 0 {
		return product, IncorrectProductID
	} else {
		product.LeagueName = product.League.GetLeagueString()
		product.RarityName = product.Rarity.GetCardRarityString()
	}

	return product, nil
}

func (p *PostgresStorage) BuyProduct(buy store.BuyProductModel) error {
	coinTr := user.CoinTransactionsModel{
		ProfileID:          buy.ProfileID,
		TransactionDetails: buy.Details,
		Amount:             buy.Coins,
		Status:             user.SuccessTransaction,
	}

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	err = p.UpdateBalance(tx, buy.ProfileID, buy.Coins)
	if err != nil {
		return err
	}
	err = p.CreateCoinTransaction(tx, coinTr)
	if err != nil {
		return err
	}

	return tx.Commit()
}
