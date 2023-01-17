package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/netwatcherio/netwatcher-control/handler"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID              primitive.ObjectID   `bson:"_id, omitempty"json:"id"`
	Email           string               `bson:"email"json:"email"` // email, will be used as username
	FirstName       string               `bson:"first_name"json:"first_name"`
	LastName        string               `bson:"last_name"json:"last_name"`
	Company         string               `bson:"company"json:"company"`
	Admin           bool                 `bson:"admin" default:"false"json:"admin"`
	Password        string               `bson:"password"json:"password"` // password in sha256?
	Sites           []primitive.ObjectID `bson:"sites"json:"sites"`       // _id's of mongo objects
	Verified        bool                 `bson:"verified"json:"sites"`    // verified, meaning email confirmation
	CreateTimestamp time.Time            `bson:"create_timestamp"json:"create_timestamp"`
}

// Create returns error if unsuccessful, error will be nil if successful
func (u *User) Create(db *mongo.Database) error {
	u.ID = primitive.NewObjectID()

	_, err := u.FromEmail(db)
	if err == nil {
		return fmt.Errorf("user exists")
	}

	mar, err := bson.Marshal(u)
	if err != nil {
		return errors.New("something went wrong marshalling user struct")
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		return errors.New("something went wrong marshalling user struct")
	}
	u.CreateTimestamp = time.Now()

	_, err = db.Collection("users").InsertOne(context.TODO(), b)
	if err != nil {
		return errors.New("something went wrong marshalling user struct")
	}

	log.Info("inserted user with the id " + u.ID.Hex())

	return nil
}

// FromID returns a user if it finds a matching user with the provided ID
func (u *User) FromID(db *mongo.Database) (*User, error) {
	var filter = bson.D{{"_id", u.ID}}
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

	if len(results) > 1 {
		return nil, errors.New("multiple users found")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	var user *User
	err = bson.Unmarshal(doc, &user)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, errors.New("something went wrong unmarshalling user data")
	}

	return user, nil
}

// FromEmail returns a user if it was able to find someone that matches the email
func (u *User) FromEmail(db *mongo.Database) (*User, error) {
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

	if len(results) > 1 {
		return nil, errors.New("multiple users found")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	var user *User
	err = bson.Unmarshal(doc, &user)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, errors.New("something went wrong unmarshalling user data")
	}

	return user, nil
}

// AddSite add a site to the user, you must create the site and ensure that the user isn't already apart of the site
// before adding the site to the user
func (u *User) AddSite(site primitive.ObjectID, db *mongo.Database) error {

	if handler.ContainsObjectID(u.Sites, site) {
		return errors.New("user is already apart of this site")
	}

	u.Sites = append(u.Sites, site)

	sites := db.Collection("users")
	_, err := sites.UpdateOne(
		context.TODO(),
		bson.M{"_id": u.ID},
		bson.D{
			{"$set", bson.D{{"sites", u.Sites}}},
		},
	)

	if err != nil {
		return errors.New("something went wrong adding site to user")
	}

	return nil
}
