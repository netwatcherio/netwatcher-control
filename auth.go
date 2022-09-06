package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/netwatcherio/netwatcher-control/control_models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func findUsers(filter bson.D, db *mongo.Database) ([]*control_models.User, error) {
	cursor, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	doc, err := bson.Marshal(results)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}

	var users []*control_models.User
	err = bson.Unmarshal(doc, &users)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	return users, nil
}

// Returns true if email exists
func getUser(email *string, db *mongo.Database) (*control_models.User, error) {
	filter := bson.D{{"email", email}}

	users, err := findUsers(filter, db)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		if len(users) > 1 {
			return nil, errors.New("2 users with same email")
		}
		return users[0], nil
	} else {
		return nil, nil
	}
}

func authUser(email *string, password string, db *mongo.Database) (*control_models.User, error) {
	user, err := getUser(email, db)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("no user exists")
	}

	if hashPassword(password) != password {
		return nil, errors.New("unable to authenticate")
	}

	return user, nil
}

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	shaHash := hex.EncodeToString(h.Sum(nil))

	return shaHash
}
