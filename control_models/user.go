package control_models

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO if generating db and stuff for first time,
// generate default admin and paste to console

type User struct {
	ID        primitive.ObjectID   `bson:"_id, omitempty"`
	Email     string               `bson:"email"` // email, will be used as username
	FirstName string               `bson:"first_name"`
	LastName  string               `bson:"last_name"`
	Admin     bool                 `bson:"admin" default:"false"`
	Password  string               `bson:"password"` // password in sha256?
	Name      string               `bson:"name"`
	Sites     []primitive.ObjectID `bson:"sites"` // _id's of mongo objects
}

// Create returns true if successful, false if otherwise with the error
func (u *User) Create(db *mongo.Database) (bool, error) {
	// todo check if already exists

	exists, err := u.UserExistsEmail(db)
	if err != nil {
		return false, err
	}

	if exists {
		return false, errors.New("user already exists")
	}

	mar, err := bson.Marshal(u)
	if err != nil {
		log.Errorf("1 %s", err)
		return false, err
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("2 %s", err)
		return false, err
	}
	result, err := db.Collection("agents").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return false, err
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return true, nil
}

// UserExistsEmail check based on wether a user with the email in user exists
func (u *User) UserExistsEmail(db *mongo.Database) (bool, error) {
	var filter = bson.D{{"email", u.Email}}

	cursor, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return false, err
	}

	//fmt.Println(results)

	if len(results) > 1 {
		return false, errors.New("multiple users match when using email")
	}

	if len(results) == 0 {
		return false, nil
	}

	if len(results) == 1 {
		return true, nil
	}

	return false, errors.New("something went wrong")
}

// UserExistsID check based on wether a user with the email in user exists
func (u *User) UserExistsID(db *mongo.Database) (bool, error) {
	var filter = bson.D{{"id", u.ID}}

	cursor, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return false, err
	}

	//fmt.Println(results)

	if len(results) > 1 {
		return false, errors.New("multiple users match when using id")
	}

	if len(results) == 0 {
		return false, nil
	}

	if len(results) == 1 {
		return true, nil
	}

	return false, errors.New("something went wrong")
}

func (u *User) GetUserFromEmail(db *mongo.Database) (*User, error) {
	var filter = bson.D{{"email", u.Email}}

	cursor, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return nil, err
	}

	var user *User
	err = bson.Unmarshal(doc, &user)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, err
	}

	return user, nil
}

func (u *User) GetUserFromID(db *mongo.Database) (*User, error) {
	var filter = bson.D{{"_id", u.ID}}

	cursor, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return nil, err
	}

	var user *User
	err = bson.Unmarshal(doc, &user)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, err
	}

	return user, nil
}
