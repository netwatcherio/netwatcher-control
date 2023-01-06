package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	MongoUri string
)

type MongoDatastore struct {
	Db      *mongo.Database
	Session *mongo.Client
	logger  *logrus.Logger
}

func NewDatastore(dbName string, logger *logrus.Logger) *MongoDatastore {

	var mongoDataStore *MongoDatastore
	db, session := Connect(dbName, logger)
	if db != nil && session != nil {

		// log statements here as well

		mongoDataStore = new(MongoDatastore)
		mongoDataStore.Db = db
		mongoDataStore.logger = logger
		mongoDataStore.Session = session
		return mongoDataStore
	}

	logger.Fatalf("Failed to Connect to database: %s", dbName)

	return nil
}

func Connect(dbName string, logger *logrus.Logger) (a *mongo.Database, b *mongo.Client) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	connectOnce.Do(func() {
		db, session = ConnectToMongo(dbName, logger)
	})

	return db, session
}

func ConnectToMongo(db string, logger *logrus.Logger) (a *mongo.Database, b *mongo.Client) {

	var err error
	session := options.Client().ApplyURI(MongoUri)
	if err != nil {
		logger.Fatal(err)
	}
	client, err := mongo.Connect(context.TODO(), session)
	if err != nil {
		logger.Fatal(err)
	}

	var DB = client.Database(db)
	logger.Infof("Successfully connected to database: %s", db)

	return DB, client
}
