package control_models

import (
	"github.com/sagostin/netwatcher-agent/agent_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	ID          primitive.ObjectID       `bson:"_id, omitempty"`
	Site        primitive.ObjectID       `bson:"site"` // _id of mongo object
	AgentConfig agent_models.AgentConfig `bson:"agent_config"`
	Pin         string                   `bson:"pin"` // used for registration & authentication
	Hash        string                   `bson:"hash"`
	// pin can be regenerated, by setting hash blank, and when registering agents, it checks for blank hashs.
}
