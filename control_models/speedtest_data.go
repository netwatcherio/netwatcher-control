package control_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/agent_models"
	"time"
)

type SpeedTestData struct {
	ID        primitive.ObjectID         `bson:"_id, omitempty"json:"id"`
	Agent     primitive.ObjectID         `bson:"agent"json:"agent"`
	Data      agent_models.SpeedTestInfo `bson:"data"json:"data"`
	Timestamp time.Time                  `bson:"timestamp"json:"timestamp"`
}
