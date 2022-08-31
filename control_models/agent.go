package control_models

import "github.com/sagostin/netwatcher-agent/agent_models"

type Agent struct {
	Site        string                   `json:"site"` // _id of mongo object
	AgentConfig agent_models.AgentConfig `json:"agent_config"`
	Pin         string                   `json:"pin"` // used for registration & authentication
	Hash        string                   `json:"hash"`
	// pin can be regenerated, by setting hash blank, and when registering agents, it checks for blank hashs.
}
