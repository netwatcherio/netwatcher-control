package handler

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiRequest struct {
	PIN   string      `json:"pin,omitempty"`
	ID    string      `json:"id,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (d *ApiRequest) GenerateCheckData(db *mongo.Database) error {
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
