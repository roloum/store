package cart

//Cart contains the information about the shopping cart and all its Items
type Cart struct {
	CartID string  `json:"cart_id"`
	Total  float32 `json:"total"`
	Count  int     `json:"count"`
	Items  []Item  `json:"items"`
}

//calculateTotal calculates the total of the shopping cart
func (c *Cart) calculateTotal() {
	for _, item := range c.Items {
		c.Total += item.Price * float32(item.Quantity)
		c.Count += item.Quantity
	}
}

//Item contains the information of an item stored in the shopping cart
type Item struct {
	ItemID      string  `json:"item_id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
}

//NewItemInfo contains the information of the new item is being added to the cart
//In case cartID is empty, a new shopping cart row is created
//Along with the item row
type NewItemInfo struct {
	CartID      string  `json:"cart_id" validate:"required"`
	ItemID      string  `json:"item_id" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required,validPrice"`
	Quantity    int     `json:"quantity" validate:"required,validQuantity"`
}

//UpdateItemInfo contains the information to update the quantity of an item in the cart
type UpdateItemInfo struct {
	CartID   string `json:"cart_id" validate:"required"`
	ItemID   string `json:"item_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,validQuantity"`
}

//DeleteItemInfo contains the information to delete an item from the cart
type DeleteItemInfo struct {
	CartID string `json:"cart_id" validate:"required"`
	ItemID string `json:"item_id" validate:"required"`
}
