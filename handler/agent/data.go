package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Data struct {
	Target    string             `json:"target,omitempty"bson:"target,omitempty"`
	ID        primitive.ObjectID `json:"id"bson:"_id"`
	CheckID   primitive.ObjectID `json:"check"bson:"check"`
	AgentID   primitive.ObjectID `json:"agent"bson:"agent"`
	Triggered bool               `json:"triggered"bson:"triggered"`
	Timestamp time.Time          `bson:"timestamp"json:"timestamp"`
	Result    interface{}        `json:"result,omitempty"bson:"result,omitempty"`
	Type      Type               `bson:"type"json:"type"`
}

func (cd *Data) Create(db *mongo.Database) error {
	// todo handle to check if agent id is set and all that... or should it be in the api section??
	cd.ID = primitive.NewObjectID()

	agentC := Check{ID: cd.CheckID}
	_, err := agentC.Get(db)
	if err != nil {
		log.Error(err)
	}

	cd.Type = agentC.Type
	crM := cd.Result.(string)
	// load types
	if agentC.Type == CtNetworkInfo {
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
	} else if agentC.Type == CtRPerf {
		var r RPerfResults
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.StopTimestamp
		cd.Result = r
	} else if agentC.Type == CtSpeedTest {
		var r SpeedTestResult
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.Timestamp
		cd.Result = r
	} else if agentC.Type == CtPing {
		var r PingResult
		err = json.Unmarshal([]byte(crM), &r)
		if err != nil {
			log.Error(err)
		}
		cd.Timestamp = r.StopTimestamp
		cd.Result = r
	}

	if (cd.Timestamp == time.Time{}) {
		// todo handle error and send alert if data that was received was not finished
		return errors.New("agent sent data with empty timestamp... skipping creation")
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

// Delete data based on provided agent ID in checkData struct
func (cd *Data) Delete(db *mongo.Database) error {
	// filter based on check ID
	var filter = bson.D{{"check", cd.CheckID}}
	if (cd.AgentID != primitive.ObjectID{}) {
		filter = bson.D{{"agent", cd.AgentID}}
	}

	_, err := db.Collection("check_data").DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}
