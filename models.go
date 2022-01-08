package main

import "time"

// Item Inventory is made of items.
type Item struct {
	ID       string `bson:"id,omitempty" json:"id,omitempty"`
	Name     string `bson:"name,omitempty" json:"name,omitempty"`
	URL      string `bson:"url,omitempty" json:"url,omitempty"`
	Quantity int64  `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Stock    int64  `bson:"stock,omitempty" json:"stock,omitempty"`
	Price    int64  `bson:"price,omitempty" json:"price,omitempty"`
}

// Iventory holds a []Item.
type Inventory struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Items []Item `json:"items,omitempty"`
}

// Shipment is its own group of items
type Shipment struct {
	ID    string
	Name  string
	Date  time.Time
	Items []Item
}

type Query interface {
	Evaluate() error
}
type InventoryQuery struct {
	StartDate time.Time
	EndDate   time.Time
}

func (inv *Inventory) Evaluate() error {
	return nil
}
