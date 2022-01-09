// Functions for Handling DB operations.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

func (db *RealDB) WriteItemToInventory(item Item) (primitive.ObjectID, error) {
	collection := db.MongoClient.Database("trade").Collection("inventory")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var id primitive.ObjectID
	defer cancel()
	res, err := collection.InsertOne(ctx, item)
	if err != nil {
		return id, err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return id, fmt.Errorf("casting error from interface{} to primitive.ObjectID")
	}
	fmt.Printf("Wrote Item #%v to DB : %v\n", res.InsertedID, res)
	fmt.Println("NewID", id)
	return id, nil
}

func (db *RealDB) ReadItems() ([]Item, error) {
	collection := db.MongoClient.Database("trade").Collection("inventory")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})
	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		fmt.Println("Collection find error")
		return nil, err
	}
	var items []Item
	err = cursor.All(context.TODO(), &items)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", items)
	return items, err
}

// EditItem edits an item with item.id with the fields.
func (db *RealDB) EditItem(item Item) (Item, error) {
	collection := db.MongoClient.Database("trade").Collection("inventory")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{"_id", item.ID}}
	res, err := collection.ReplaceOne(ctx, filter, item)
	if err != nil {
		return item, err
	}
	fmt.Printf("%+v\n", res)
	return item, err
}

func (db *RealDB) DeleteItem(key string) error {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return err
	}
	collection := db.MongoClient.Database("trade").Collection("inventory")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	fmt.Println("Deleted count: ", res.DeletedCount)
	if res.DeletedCount == 1 {
		return nil
	}
	return fmt.Errorf("nothing deleted from database")
}
