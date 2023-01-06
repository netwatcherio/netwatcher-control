package handler

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TODO if generating db and stuff for first time,
// generate default admin and paste to console

type User struct {
	ID              primitive.ObjectID   `bson:"_id, omitempty"json:"id"`
	Email           string               `bson:"email"json:"email"` // email, will be used as username
	FirstName       string               `bson:"first_name"json:"first_name"`
	LastName        string               `bson:"last_name"json:"last_name"`
	Admin           bool                 `bson:"admin" default:"false"json:"admin"`
	Password        string               `bson:"password"`             // password in sha256?
	Sites           []primitive.ObjectID `bson:"sites"json:"sites"`    // _id's of mongo objects
	Verified        bool                 `bson:"verified"json:"sites"` // verified, meaning email confirmation
	CreateTimestamp time.Time            `bson:"create_timestamp"json:"create_timestamp"`
}

type RegisterUser struct {
	Email           string `bson:"email"form:"email"` // email, will be used as username
	FirstName       string `bson:"first_name"form:"first_name"`
	LastName        string `bson:"last_name"form:"last_name"`
	Password        string `bson:"password"form:"password"`                 // password in sha256?
	PasswordConfirm string `bson:"password_confirm"form:"password_confirm"` // password in sha256?
}

type LoginUser struct {
	Email    string `bson:"email"form:"email"`       // email, will be used as username
	Password string `bson:"password"form:"password"` // password in sha256?
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
		return false, errors.New("something went wrong")
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("2 %s", err)
		return false, errors.New("something went wrong")
	}
	u.CreateTimestamp = time.Now()

	result, err := db.Collection("users").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return false, errors.New("something went wrong")
	}

	fmt.Printf(" with _id: %v\n", result.InsertedID)
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
		return false, errors.New("something went wrong")
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
	var filter = bson.D{{"_id", u.ID}}

	cursor, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return false, errors.New("something went wrong")
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

	if len(results) < 1 {
		return nil, errors.New("no user found")
	}
	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return nil, errors.New("something went wrong")
	}

	var user *User
	err = bson.Unmarshal(doc, &user)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, errors.New("something went wrong")
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
		return nil, errors.New("something went wrong")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return nil, errors.New("something went wrong")
	}

	var user *User
	err = bson.Unmarshal(doc, &user)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, errors.New("something went wrong")
	}

	return user, nil
}

// AddSite add the id of a site
func (u *User) AddSite(site primitive.ObjectID, db *mongo.Database) (bool, error) {
	// add member with the provided role

	u.Sites = append(u.Sites, site)

	/*sitB, err := bson.Marshal(u.Sites)
	if err != nil {
		log.Errorf("69 %s", err)
	}*/

	sites := db.Collection("users")
	result, err := sites.UpdateOne(
		context.TODO(),
		bson.M{"_id": u.ID},
		bson.D{
			{"$set", bson.D{{"sites", u.Sites}}},
		},
	)
	log.Warnf("%s", result)

	if err != nil {
		log.Error(err)
		return false, err
	}

	return true, nil
}

// todo when deleting site remove from user document as well
