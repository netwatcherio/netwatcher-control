package control_models

import (
	"github.com/netwatcherio/netwatcher-agent/agent_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MtrData struct {
	ID        primitive.ObjectID       `bson:"_id, omitempty"`
	Agent     primitive.ObjectID       `bson:"agent"`
	Data      []agent_models.MtrTarget `bson:"data"`
	Timestamp time.Time                `bson:"timestamp"`
}
