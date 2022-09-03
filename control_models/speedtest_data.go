package control_models

import (
	"github.com/sagostin/netwatcher-agent/agent_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SpeedTestData struct {
	ID        primitive.ObjectID         `bson:"_id, omitempty"`
	Agent     primitive.ObjectID         `bson:"agent"`
	Data      agent_models.SpeedTestInfo `bson:"data"`
	Timestamp time.Time                  `bson:"timestamp"`
}
