package control_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/agent_models"
	"time"
)

type IcmpData struct {
	ID        primitive.ObjectID        `bson:"_id, omitempty"json:"id"`
	Agent     primitive.ObjectID        `bson:"agent"json:"agent"` // _id of mongo object
	Data      []agent_models.IcmpTarget `bson:"data"json:"data"`
	Timestamp time.Time                 `bson:"timestamp"json:"timestamp"`
}
