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
func TestCreateAddItem(t *testing.T) {

	handler, _ := New(&test.MockDynamoDB{}, StoreTable)
	ni := NewItem{ItemID: "123", Description: "desc", Price: 1, Quantity: 1}
	_, err := handler.CreateAndAddItem(context.Background(), ni)
	if err != nil {
		t.Fatalf("Received error while creating shopping cart: %s", err.Error())
	}

}
