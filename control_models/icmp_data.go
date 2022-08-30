package control_models

import "github.com/sagostin/netwatcher-agent/agent_models"

type IcmpData struct {
	Agent string                  `json:"agent"` // _id of mongo object
	Data  agent_models.IcmpTarget `json:"data"`
}
