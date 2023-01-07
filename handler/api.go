package handler

import (
	"github.com/netwatcherio/netwatcher-agent/checks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Data struct {
	PIN    string             `json:"pin,omitempty"`
	ID     string             `json:"id,omitempty"`
	Checks []checks.CheckData `json:"checks,omitempty"`
	Error  string             `json:"error,omitempty"`
}

func (d *Data) GenerateCheckData(db *mongo.Database) error {
	hexId, err := primitive.ObjectIDFromHex(d.ID)
	if err != nil {
		return err
	}

	agentCheck := AgentCheck{AgentID: hexId}
	_, err = agentCheck.Get(db)
	if err != nil {
		return err
	}
	return nil
}
