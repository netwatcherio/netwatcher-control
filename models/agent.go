package models

import (
	"github.com/netwatcherio/netwatcher-agent/checks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Agent struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"json:"id"`
	Name      string             `bson:"name"json:"name"`
	Site      primitive.ObjectID `bson:"site"json:"site"` // _id of mongo object
	Pin       string             `bson:"pin"json:"pin"`   // used for registration & authentication
	Heartbeat time.Time          `bson:"heartbeat,omitempty"json:"heartbeat,omitempty"`
	CheckData []checks.CheckData `json:"check_data"bson:"check_data"`
	// pin can be regenerated, by setting hash blank, and when registering agents, it checks for blank hashs.
}
