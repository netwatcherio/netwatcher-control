package handler

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CheckData struct {
	Target    string             `json:"target,omitempty"bson:"target,omitempty"`
	ID        primitive.ObjectID `json:"id"bson:"_id"`
	CheckID   primitive.ObjectID `json:"check"bson:"check"`
	AgentID   primitive.ObjectID `json:"agent"bson:"agent"`
	Triggered bool               `json:"triggered"bson:"triggered,omitempty"`
	Timestamp time.Time          `bson:"timestamp"json:"timestamp"`
	Result    interface{}        `json:"result"bson:"result,omitempty"`
	Type      CheckType          `bson:"type"json:"type"`
}

func (cd *CheckData) Create(db *mongo.Database) error {
	// todo handle to check if agent id is set and all that... or should it be in the api section??
	cd.ID = primitive.NewObjectID()

	agentC := AgentCheck{ID: cd.CheckID}
	_, err := agentC.Get(db)
	if err != nil {
		log.Error(err)
	}

	cd.Type = agentC.Type

	crM := cd.Result.(string)

	// load types
	if agentC.Type == CtNetinfo {
		var r NetResult
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.Timestamp
		cd.Result = r
	} else if agentC.Type == CtMtr {
		var r MtrResult
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.StopTimestamp
		cd.Result = r
	} else if agentC.Type == CtRperf {
		var r RPerfResults
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.StopTimestamp
		cd.Result = r
	} else if agentC.Type == CtSpeedtest {
		var r SpeedTest
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.Timestamp
		cd.Result = r
	}

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

	cAgent := Agent{ID: cd.AgentID}
	err = cAgent.UpdateHeartbeat(db)
	if err != nil {
		return err
	}

	fmt.Printf("created check data with id: %v\n", result.InsertedID)
	return nil
}
