package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	mongoUri string
)

type MongoDatastore struct {
	db      *mongo.Database
	Session *mongo.Client
	logger  *logrus.Logger
}

func NewDatastore(dbName string, logger *logrus.Logger) *MongoDatastore {

	var mongoDataStore *MongoDatastore
	db, session := connect(dbName, logger)
	if db != nil && session != nil {

		// log statements here as well

		mongoDataStore = new(MongoDatastore)
		mongoDataStore.db = db
		mongoDataStore.logger = logger
		mongoDataStore.Session = session
		return mongoDataStore
	}

	logger.Fatalf("Failed to connect to database: %s", dbName)

	return nil
}

func connect(dbName string, logger *logrus.Logger) (a *mongo.Database, b *mongo.Client) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	connectOnce.Do(func() {
		db, session = connectToMongo(dbName, logger)
	})

	return db, session
}

func connectToMongo(db string, logger *logrus.Logger) (a *mongo.Database, b *mongo.Client) {

	var err error
	session := options.Client().ApplyURI(mongoUri)
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
