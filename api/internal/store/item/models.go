package item

//Item contains the information of an item
type Item struct {
	ItemID      string  `json:"item_id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

//List contains a list of items
type List struct {
	Items []Item `json:"items"`
}
