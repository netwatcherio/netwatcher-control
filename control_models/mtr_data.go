package control_models

import "github.com/sagostin/netwatcher-agent/agent_models"

type MtrData struct {
	Agent string `json:"agent"`
	Data  agent_models.MtrTarget
}
