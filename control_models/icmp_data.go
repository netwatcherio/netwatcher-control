package control_models

import (
	"github.com/netwatcherio/netwatcher-agent/agent_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IcmpData struct {
	ID        primitive.ObjectID        `bson:"_id, omitempty"`
	Agent     primitive.ObjectID        `bson:"agent"` // _id of mongo object
	Data      []agent_models.IcmpTarget `bson:"data"`
	Timestamp time.Time                 `bson:"timestamp"`
}
