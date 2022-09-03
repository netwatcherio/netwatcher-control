package control_models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TODO if generating db and stuff for first time,
// generate default admin and paste to console

type User struct {
	ID        primitive.ObjectID   `bson:"_id, omitempty"`
	Email     string               `json:"email"` // username
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Admin     bool                 `json:"admin"`
	Password  string               `json:"password"` // password in sha256?
	Name      string               `json:"name"`
	Sites     []primitive.ObjectID `bson:"sites"` // _id's of mongo objects
}
