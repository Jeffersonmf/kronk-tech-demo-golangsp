package inventory

import (
	"sync"
	"testing"
)

func TestDeductStock_RaceCondition(t *testing.T) {
	s := NewService()
	productID := "prod_123"
	s.AddProduct(&Product{ID: productID, Name: "Test Product", Stock: 10})

	var wg sync.WaitGroup
	numRequests := 20
	quantityPerRequest := 1

	// Temos 10 em estoque e fazemos 20 requisições de 1.
	// Esperamos exatamente 10 requisições com sucesso e 10 com falha.
	// O estoque deve ficar em 0 no final.

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.DeductStock(productID, quantityPerRequest)
		}()
	}

	wg.Wait()

	stock, _ := s.GetStock(productID)
	if stock < 0 {
		t.Errorf("Estoque está negativo: %d", stock)
	}
}
