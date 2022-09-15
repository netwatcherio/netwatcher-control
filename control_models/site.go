package control_models

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Site struct {
	ID      primitive.ObjectID `bson:"_id, omitempty"`
	Name    string             `bson:"name"`
	Members []SiteMember       `bson:"members"`
}

type SiteMember struct {
	User primitive.ObjectID `bson:"user"`
	Role int                `bson:"role"`
	// roles: 0=READ ONLY, 1=READ-WRITE (Create only), 2=ADMIN (Delete Agents), 3=OWNER (Delete Sites)
	// ADMINS can regenerate agent pins
}

func (s *Site) CreateSite(name string, owner primitive.ObjectID, db *mongo.Database) (bool, error) {
	mar, err := bson.Marshal(s)
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
	_, err = db.Collection("sites").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return false, err
	}
	return true, nil
}

// IsMember check if a user id is a member in the site
func (s *Site) IsMember(id primitive.ObjectID) bool {
	// check if the site contains the member with the provided id
	for _, m := range s.Members {
		if m.User == id {
			return true
		}
	}

	return false
}

// AddMember Add a member to the site then update document
func (s *Site) AddMember(id primitive.ObjectID, db *mongo.Database) (bool, error) {
	// add member with the provided role
	if s.IsMember(id) {
		return false, errors.New("already a member")
	}

	newMember := SiteMember{
		User: id,
		Role: 1,
	}

	s.Members = append(s.Members, newMember)

	sites := db.Collection("sites")
	_, err := sites.UpdateOne(
		context.TODO(),
		bson.M{"_id": s.ID},
		bson.D{
			{"$set", bson.D{{"members", s.Members}}},
		},
	)
	if err != nil {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (s *Site) SiteExists() {

}

// todo handle if site already exists, or already has the member in it
