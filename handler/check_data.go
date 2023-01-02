package handler

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CheckData struct {
	Target    string             `json:"address,omitempty"bson:"target,omitempty"`
	ID        primitive.ObjectID `json:"id"bson:"_id"`
	CheckID   primitive.ObjectID `json:"check"bson:"check"`
	AgentID   primitive.ObjectID `json:"agent"bson:"agent"`
	Triggered bool               `json:"triggered"bson:"triggered,omitempty"`
	Result    interface{}        `json:"result"bson:"result,omitempty"`
}

func (cd *CheckData) Create(db *mongo.Database) error {
	// todo handle to check if agent id is set and all that... or should it be in the api section??
	cd.ID = primitive.NewObjectID()

	mar, err := bson.Marshal(cd)
	if err != nil {
		log.Errorf("error marshalling check data when creating: %s", err)
		return err
	}

	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("error unmarhsalling check data when creating: %s", err)
		return err
	}
	result, err := db.Collection("check_data").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("error inserting to database: %s", err)
		return err
	}

	fmt.Printf("created check data with id: %v\n", result.InsertedID)
	return nil
}
