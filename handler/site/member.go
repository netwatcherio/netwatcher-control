package site

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Role string

const (
	SrREAD_ONLY  Role = "READ_ONLY"
	SrREAD_WRITE Role = "READ_WRITE"
	SrADMIN      Role = "ADMIN"
	SrOWNER      Role = "OWNER"
)

type Member struct {
	User primitive.ObjectID `bson:"user"json:"user"`
	Role Role               `bson:"role"json:"role"`
	// roles: 0=READ ONLY, 1=READ-WRITE (Create only), 2=ADMIN (Delete Agents), 3=OWNER (Delete Sites)
	// ADMINS can regenerate agent pins
}

type AddSiteMember struct {
	Email string `json:"email"form:"email"`
	Role  Role   `json:"role"form:"role"`
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
func (s *Site) AddMember(id primitive.ObjectID, role Role, db *mongo.Database) error {
	// add member with the provided role
	if s.IsMember(id) {
		return errors.New("already a member")
	}

	newMember := Member{
		User: id,
		Role: role,
	}

	s.Members = append(s.Members, newMember)
	j, _ := json.Marshal(s.Members)
	log.Warnf("%s", j)

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
		return err
	}
	return nil
}
