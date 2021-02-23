package cart

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/roloum/store/api/internal/test"
	"github.com/rs/zerolog"
)

const (
	StoreTable = "Store"
)

type (
	cartTest struct {
		desc string
		item *NewItemInfo
		err  error
	}
)

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	test.SetEnvironment()
}

//TestCreateCart tests the CreateAddItem method that creates a new shopping cart
func TestCreateAddItem(t *testing.T) {

	handler, _ := New(&test.MockDynamoDB{}, StoreTable)

	tests := []cartTest{
		{
			desc: ErrCreateCartWithExistingCartID.Error(),
			item: &NewItemInfo{
				CartID:      "wrongID",
				ItemID:      "",
				Description: "",
				Price:       0,
				Quantity:    0,
			},
			err: ErrCreateCartWithExistingCartID,
		},
	}
	tests = append(tests, getSuccessAddItem())
	tests = append(tests, getAddItemTestCases()...)
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := handler.CreateAndAddItem(context.Background(), tc.item)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Expected: %v. Received: %v", tc.err, err)
			}
		})
	}

}

//TestAddItem tests the AddItem to an existing shopping cart
func TestAddItem(t *testing.T) {
	handler, _ := New(&test.MockDynamoDB{}, StoreTable)

	//Test AddItem without cartID
	t.Run(ErrCartIDIsEmpty, func(t *testing.T) {
		expectedErr := errors.New(ErrCartIDIsEmpty)
		_, err := handler.AddItem(context.Background(), &NewItemInfo{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Errorf("Expected: %v. Received: %v", expectedErr, err)
		}
	})

	//Create cart and retrieve the cartID
	var cartID string
	newCartItem := getSuccessAddItem()
	t.Run(newCartItem.desc, func(t *testing.T) {

		c, err := handler.CreateAndAddItem(context.Background(), newCartItem.item)
		if !reflect.DeepEqual(err, newCartItem.err) {
			t.Errorf("Expected: %v. Received: %v", newCartItem.err, err)
		}

		cartID = c.CartID
	})

	t.Logf("CartID: %s", cartID)

	//Add cartID to the test cases
	tests := getAddItemTestCases()
	for _, test := range tests {
		test.item.CartID = cartID
	}

	//Test adding items for the newly created shopping cart
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := handler.AddItem(context.Background(), tc.item)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Expected: %v. Received: %v", tc.err, err)
			}
		})
	}
}

//getSuccessCartItem returns a successful test case that creates a shopping cart
func getSuccessAddItem() cartTest {
	return cartTest{
		desc: "Success",
		item: &NewItemInfo{
			ItemID:      "11aa",
			Description: "Some item description",
			Price:       1,
			Quantity:    1,
		},
		err: nil,
	}
}

//getAddItemTestCases returns the test cases that are going to be used for
//testing both TestCreateAddItem and TestAddItem
func getAddItemTestCases() []cartTest {
	return []cartTest{
		{
			desc: ErrItemIDIsEmpty,
			item: &NewItemInfo{
				ItemID:      "",
				Description: "",
				Price:       0,
				Quantity:    0,
			},
			err: errors.New(ErrItemIDIsEmpty),
		},
		{
			desc: ErrDescriptionIsEmpty,
			item: &NewItemInfo{
				ItemID:      "11aa",
				Description: "",
				Price:       0,
				Quantity:    0,
			},
			err: errors.New(ErrDescriptionIsEmpty),
		},
		{
			desc: ErrPriceIsEmpty,
			item: &NewItemInfo{
				ItemID:      "11aa",
				Description: "Some item description",
				Price:       0,
				Quantity:    0,
			},
			err: errors.New(ErrPriceIsEmpty),
		},
		{
			desc: ErrPriceIsInvalid,
			item: &NewItemInfo{
				ItemID:      "11aa",
				Description: "Some item description",
				Price:       -1,
				Quantity:    0,
			},
			err: errors.New(ErrPriceIsInvalid),
		},
		{
			desc: ErrQuantityIsEmpty,
			item: &NewItemInfo{
				ItemID:      "11aa",
				Description: "Some item description",
				Price:       1,
				Quantity:    0,
			},
			err: errors.New(ErrQuantityIsEmpty),
		},
		{
			desc: ErrQuantityIsInvalid,
			item: &NewItemInfo{
				ItemID:      "11aa",
				Description: "Some item description",
				Price:       1,
				Quantity:    -1,
			},
			err: errors.New(ErrQuantityIsInvalid),
		},
		{
			desc: "AddItemSuccess",
			item: &NewItemInfo{
				ItemID:      "22bb",
				Description: "Some item description",
				Price:       1,
				Quantity:    1,
			},
			err: nil,
		},
	}
}
