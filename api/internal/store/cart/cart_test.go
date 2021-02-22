package cart

import (
	"context"
	"testing"

	"github.com/roloum/store/api/internal/test"
	"github.com/rs/zerolog"
)

const (
	StoreTable = "Store"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	test.SetEnvironment()
}

//TestCreateCart tests the Create method that creates a new shopping cart
func TestCreateCart(t *testing.T) {

	cart, _ := New(&test.MockDynamoDB{}, StoreTable)
	_, err := cart.Create(context.Background())
	if err != nil {
		t.Fatalf("Received error while creating shopping cart: %s", err.Error())
	}

}
