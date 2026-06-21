package inventory

import (
	"errors"
	"fmt"
)

var ErrInsufficientStock = errors.New("estoque insuficiente")

type Product struct {
	ID    string
	Name  string
	Stock int
}

type Service struct {
	products map[string]*Product
}

func NewService() *Service {
	return &Service{
		products: make(map[string]*Product),
	}
}

func (s *Service) AddProduct(p *Product) {
	s.products[p.ID] = p
}

// DeductStock é a função com o bug. Ela tem uma race condition.
func (s *Service) DeductStock(productID string, quantity int) error {
	p, ok := s.products[productID]
	if !ok {
		return fmt.Errorf("produto %s não encontrado", productID)
	}

	// BUG: Race condition aqui. Múltiplas goroutines podem passar por essa
	// verificação antes que a subtração aconteça.
	if p.Stock < quantity {
		return ErrInsufficientStock
	}

	// Simula algum trabalho/latência para tornar a race mais visível
	// (não é estritamente necessário, mas ajuda na demonstração).
	p.Stock -= quantity

	return nil
}

func (s *Service) GetStock(productID string) (int, error) {
	p, ok := s.products[productID]
	if !ok {
		return 0, fmt.Errorf("produto %s não encontrado", productID)
	}
	return p.Stock, nil
}
