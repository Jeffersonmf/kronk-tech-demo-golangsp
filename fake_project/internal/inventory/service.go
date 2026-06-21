package inventory

import (
	"errors"
	"fmt"
	"sync"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type Product struct {
	ID    string
	Name  string
	Stock int
}

type Service struct {
	products map[string]*Product
	mu       sync.Mutex
}

func NewService() *Service {
	return &Service{
		products: make(map[string]*Product),
	}
}

func (s *Service) AddProduct(p *Product) {
	s.products[p.ID] = p
}

// DeductStock deducts the stock of a product.
func (s *Service) DeductStock(productID string, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[productID]
	if !ok {
		return fmt.Errorf("product %s not found", productID)
	}

	if p.Stock < quantity {
		return ErrInsufficientStock
	}

	p.Stock -= quantity

	return nil
}

func (s *Service) GetStock(productID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[productID]
	if !ok {
		return 0, fmt.Errorf("product %s not found", productID)
	}
	return p.Stock, nil
}
