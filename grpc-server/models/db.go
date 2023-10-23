package models

import (
	"context"

	"github.com/aeon/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConn is a struct representing the MongoDB connection
type DBConn struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

type DBHandlers interface{}

// NewMongoDBConnection creates a new MongoDB connection
func NewDBConnection(ctx context.Context, collectionName string) (*DBConn, error) {
	dbEnv := utils.LoadEnvConfig("DB")
	dbName := dbEnv.GetString("NAME")
	dbURI := dbEnv.GetString("URI")

	clientOptions := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Errorf("Error connecting to the database %s", err.Error())
		return nil, err
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &DBConn{
		Client:     client,
		Database:   db,
		Collection: collection,
	}, nil
}

func (conn *DBConn) InsertRecord(ctx context.Context, record interface{}) error {
	collection := conn.Collection
	_, err := collection.InsertOne(ctx, record)
	if err != nil {
		logrus.Errorf("Error inserting new record %v in to db: %s", record, err.Error())
		return err
	}
	return nil
}

func (conn *DBConn) UpdateRecord(ctx context.Context, filter interface{}, record interface{}) error {
	collection := conn.Collection
	_, err := collection.UpdateOne(ctx, filter, record)
	if err != nil {
		logrus.Errorf("Error updating a record in to db: %s", err.Error())
		return err
	}
	return nil
}

func (conn *DBConn) DeleteRecord(ctx context.Context, ID string) error {
	collection := conn.Collection
	_, err := collection.DeleteOne(ctx, bson.M{"uuid": ID})
	if err != nil {
		logrus.Errorf("Error deleting a book with uuid: %s, err: %s", ID, err.Error())
		return err
	}
	return nil
}

func (conn *DBConn) FindAllRecords(ctx context.Context, records interface{}, filter interface{}) error {
	collection := conn.Collection
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		logrus.Errorf("Error finding books: %s", err.Error())
		return err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, records); err != nil {
		logrus.Errorf("Error iterating through the cursor: %s", err.Error())
		return err
	}
	return nil
}
