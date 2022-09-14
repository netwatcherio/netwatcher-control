package control_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Site struct {
	ID      primitive.ObjectID `bson:"_id, omitempty"`
	Name    string             `bson:"name"`
	Members []struct {
		User primitive.ObjectID `bson:"user"` // _id
		Role int                `bson:"role"`
		// roles: 0=READ ONLY, 1=READ-WRITE (Create only), 2=ADMIN (Delete Agents), 3=OWNER (Delete Sites)
		// ADMINS can regenerate agent pins
	}
}

func (s *Site) AddMember(id primitive.ObjectID, db *mongo.Database) (bool, error) {

	return false, nil
}

// CreateSite create site with the owner's ID as one of the members with the role of 3
func (s *Site) CreateSite(id primitive.ObjectID, db *mongo.Database) (bool, error) {

	return false, nil
}

func (s *Site) SiteExists() {

}

// todo handle if site already exists, or already has the member in it
