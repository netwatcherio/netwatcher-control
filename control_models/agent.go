package control_models

import (
	"github.com/netwatcherio/netwatcher-agent/agent_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Agent struct {
	ID          primitive.ObjectID       `bson:"_id, omitempty"json:"id"`
	Name        string                   `bson:"name"json:"name"`
	Site        primitive.ObjectID       `bson:"site"json:"site"` // _id of mongo object
	AgentConfig agent_models.AgentConfig `bson:"agent_config"json:"agent_config"`
	Pin         string                   `bson:"pin"json:"pin"` // used for registration & authentication
	Hash        string                   `bson:"hash"json:"hash"`
	Heartbeat   time.Time                `bson:"heartbeat"json:"heartbeat"`
	// pin can be regenerated, by setting hash blank, and when registering agents, it checks for blank hashs.
}

type CreateAgent struct {
	Name        string `json:"name"form:"name"`
	IcmpTargets string `json:"icmpTargets"form:"icmpTargets"`
	MtrTargets  string `json:"mtrTargets"form:"mtrTargets"`
}

type AgentStats struct {
	ID          primitive.ObjectID       `json:"id"`
	Name        string                   `json:"name"`
	Heartbeat   time.Time                `json:"heartbeat"`
	NetworkInfo agent_models.NetworkInfo `json:"network_info"`
	LastSeen    time.Duration            `json:"last_seen"`
	Online      bool                     `json:"online"`
}

type AgentStatsList struct {
	List []AgentStats `json:"list"`
}
