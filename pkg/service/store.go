package service

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"log"
)

func NewStoreService(storage StoreStorage) *StoreService {
	return &StoreService{
		storage: storage,
	}
}

type StoreStorage interface {
	GetAllProducts() ([]store.Product, error)
	GetProductByID(id int) (store.Product, error)
	BuyProduct(buy store.BuyProductModel) error
}

type StoreService struct {
	storage StoreStorage
}

func (s *StoreService) GetAllProducts() ([]store.Product, error) {

	products, err := s.storage.GetAllProducts()
	if err != nil {
		log.Println("Service. GetAllProducts:", err)
		return products, err
	}

	return products, nil
}

func (s *StoreService) BuyProduct(buy store.BuyProductModel) error {

	product, err := s.storage.GetProductByID(buy.ID)
	if err != nil {
		log.Println("Service. GetProductByID:", err)
		return err
	}
	buy.Coins = -product.Price
	buy.Details = "Покупка: " + product.ProductName

	err = s.storage.BuyProduct(buy)
	if err != nil {
		log.Println("Service. BuyProduct:", err)
		return err
	}

	return nil
}
