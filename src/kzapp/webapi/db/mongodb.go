package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	log.Printf("正在連接到 MongoDB: %s", uri)

	clientOpts := createClientOptions(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect from MongoDB: %v", err)
	}
	// 驗證連接
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("成功連接到 MongoDB")
	return client
}

func createClientOptions(uri string) *options.ClientOptions {
	// Specify BSON options that cause the driver to fallback to "json"
	// struct tags if "bson" struct tags are missing, marshal nil Go maps as
	// empty BSON documents, and marshals nil Go slices as empty BSON
	// arrays.
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}

	clientOpts := options.Client().
		ApplyURI(uri).
		SetBSONOptions(bsonOpts).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(10 * time.Second).
		SetConnectTimeout(5 * time.Second)
	return clientOpts
}

func NewFirstMongoDB(uri string) *FirstDB {
	client := InitMongoDB(uri)
	userDao := &UserMongoDaoImpl{db: client.Database("kzapp")}
	return &FirstDB{db: client, UserDao: userDao}
}
