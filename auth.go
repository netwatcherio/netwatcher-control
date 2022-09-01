package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func Auth(email *string, password *string, db *mongo.Database) {
	filter := bson.D{{"email", email}}

	cursor, err := db.Collection("").Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	shaHash := hex.EncodeToString(h.Sum(nil))

	return shaHash
}
