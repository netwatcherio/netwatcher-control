package control_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/agent_models"
	"time"
)

type NetworkData struct {
	ID        primitive.ObjectID       `bson:"_id, omitempty"json:"id"`
	Agent     primitive.ObjectID       `bson:"agent"json:"agent"`
	Data      agent_models.NetworkInfo `bson:"data"json:"data"`
	Timestamp time.Time                `bson:"timestamp"json:"timestamp"`
}
