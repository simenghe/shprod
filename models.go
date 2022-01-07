package main

// Item Inventory is made of items.
type Item struct {
	Key      string
	Name     string
	Quantity int64
}

// Shipment is its own group of items
type Shipment struct {
	Items []Item
}
