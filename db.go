// Functions for Handling DB operations.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Currently just operating one inventory.
type FakeDB struct {
	FakeInventory Inventory
}

type RealDB struct {
	MongoClient *mongo.Client
}

func NewDB(uri string) *RealDB {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	// Also do a ping when connecting.
	return &RealDB{
		MongoClient: client,
	}
}

func (db *RealDB) Disconnect() error {
	return db.MongoClient.Disconnect(context.Background())
}

func (db *RealDB) WriteItemToInventory(item Item) error {
	collection := db.MongoClient.Database("trade").Collection("inventory")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, item)
	fmt.Printf("Wrote Item #%v to DB : %v\n", res.InsertedID, res)
	return err
}

func NewFakeDB() *FakeDB {
	db := &FakeDB{
		FakeInventory: Inventory{
			ID:    "0",
			Name:  "FakeInventory",
			Items: make([]Item, 0),
		},
	}
	db.FakeInventory.Items = append(db.FakeInventory.Items, Item{Name: "Test0", ID: "0"})
	db.FakeInventory.Items = append(db.FakeInventory.Items, Item{Name: "Test1", ID: "1"})
	return db
}

func (db *FakeDB) WriteItemToInventory(item Item) error {
	db.FakeInventory.Items = append(db.FakeInventory.Items, item)
	return nil
}

// EditItem changes the Item with ID, and changes the fields to the newItem.
func (db *FakeDB) EditItem(ID string, newItem Item) error {
	for i, item := range db.FakeInventory.Items {
		if item.ID == ID {
			// Copy all the fields that were modified, keep the ID the same.
			db.FakeInventory.Items[i] = newItem
			db.FakeInventory.Items[i].ID = ID
			return nil
		}
	}
	return fmt.Errorf("could not find Item with ID : %v", ID)
}
