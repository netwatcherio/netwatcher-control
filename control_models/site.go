package control_models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Site struct {
	ID      primitive.ObjectID   `bson:"_id, omitempty"`
	Agents  []primitive.ObjectID `bson:"agents"` // _id of mongo objects
	Members []struct {
		User primitive.ObjectID `bson:"user"` // _id
		Role int                `json:"role"`
		// roles: 0=READ ONLY, 1=READ-WRITE (Create only), 2=ADMIN (Delete Agents), 3=OWNER (Delete Sites)
		// ADMINS can regenerate agent pins
	}
}
