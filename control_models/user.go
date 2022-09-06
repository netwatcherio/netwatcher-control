package control_models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TODO if generating db and stuff for first time,
// generate default admin and paste to console

type User struct {
	ID        primitive.ObjectID   `bson:"_id, omitempty"`
	Email     string               `bson:"email"` // username
	FirstName string               `bson:"first_name"`
	LastName  string               `bson:"last_name"`
	Admin     bool                 `bson:"admin"`
	Password  string               `bson:"password"` // password in sha256?
	Name      string               `bson:"name"`
	Sites     []primitive.ObjectID `bson:"sites"` // _id's of mongo objects
}
