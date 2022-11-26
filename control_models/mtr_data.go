package control_models

import (
	"github.com/netwatcherio/netwatcher-agent/agent_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MtrData struct {
	ID        primitive.ObjectID       `bson:"_id, omitempty"json:"id"`
	Agent     primitive.ObjectID       `bson:"agent"json:"agent"`
	Data      []agent_models.MtrTarget `bson:"data"json:"data"`
	Timestamp time.Time                `bson:"timestamp"json:"timestamp"`
}
