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

	// We have 10 in stock, we make 20 requests of 1.
	// We expect exactly 10 requests to succeed and 10 to fail.
	// Stock should be 0 at the end.

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
		t.Errorf("Stock is negative: %d", stock)
	}
}
