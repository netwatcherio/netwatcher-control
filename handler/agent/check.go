package agent

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Check struct {
	Type            Type               `json:"type"bson:"type"`
	Target          string             `json:"target"bson:"target"`
	ID              primitive.ObjectID `json:"id"bson:"_id"`
	AgentID         primitive.ObjectID `json:"agent"bson:"agent"`
	Duration        int                `json:"duration'"bson:"duration"`
	Count           int                `json:"count"bson:"count"`
	Triggered       bool               `json:"triggered"bson:"triggered"`
	Server          bool               `json:"server"bson:"server"`
	Pending         bool               `json:"pending"bson:"pending"`
	Interval        int                `json:"interval"bson:"interval"`
	CreateTimestamp time.Time          `bson:"create_timestamp"json:"create_timestamp"`
}

type CheckRequest struct {
	Limit          int64     `json:"limit"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp"`
	Recent         bool      `json:"recent"`
}

func (ac *Check) GetData(req CheckRequest, db *mongo.Database) ([]*Data, error) {
	opts := options.Find().SetLimit(req.Limit)

	var filter = bson.D{{"check", ac.ID}}
	if ac.AgentID != (primitive.ObjectID{0}) {
		filter = bson.D{{"agent", ac.AgentID}, {"type", ac.Type}}
	}

	var timeFilter bson.M

	if req.Recent {
		opts = opts.SetSort(bson.D{{"timestamp", -1}})
	} else {
		if ac.Type == CtRPerf {
			timeFilter = bson.M{
				"check": ac.ID,
				"result.stop_timestamp": bson.M{
					"$gt": req.StartTimestamp,
					"$lt": req.EndTimestamp,
				}}
		} else {
			timeFilter = bson.M{
				"check": ac.ID,
				"timestamp": bson.M{
					"$gt": req.StartTimestamp,
					"$lt": req.EndTimestamp,
				}}
		}
	}

	var cursor *mongo.Cursor
	var err error

	if req.Recent {
		cursor, err = db.Collection("check_data").Find(context.TODO(), filter, opts)
		if err != nil {
			return nil, err
		}
	} else {
		cursor, err = db.Collection("check_data").Find(context.TODO(), timeFilter, opts)
		if err != nil {
			return nil, err
		}
	}
	var results []bson.D
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) <= 0 {
		return nil, errors.New("no data matches the provided check id")
	}

	if len(results) <= 0 {
		return nil, errors.New("no data found")
	}

	var checkData []*Data

	for _, r := range results {
		var cData Data
		doc, err := bson.Marshal(r)
		if err != nil {
			log.Errorf("1 %s", err)
			return nil, err
		}
		err = bson.Unmarshal(doc, &cData)
		if err != nil {
			log.Errorf("22 %s", err)
			return nil, err
		}

		checkData = append(checkData, &cData)
	}

	return checkData, nil
}

func (ac *Check) Create(db *mongo.Database) error {
	ac.ID = primitive.NewObjectID()

	mar, err := bson.Marshal(ac)
	if err != nil {
		log.Errorf("error marshalling agent check when creating: %s", err)
		return err
	}

	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("error unmarhsalling agent check when creating: %s", err)
		return err
	}
	result, err := db.Collection("agent_check").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("error inserting to database: %s", err)
		return err
	}

	fmt.Printf("created agent check with id: %v\n", result.InsertedID)
	return nil
}

func (ac *Check) Get(db *mongo.Database) ([]*Check, error) {
	var filter = bson.D{{"_id", ac.ID}}

	if ac.AgentID != (primitive.ObjectID{0}) {
		filter = bson.D{{"agent", ac.AgentID}}
	}

	cursor, err := db.Collection("agent_check").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	//fmt.Println(results)

	if ac.AgentID == (primitive.ObjectID{0}) {
		if len(results) > 1 {
			return nil, errors.New("multiple sites match when using id")
		}

		if len(results) == 0 {
			return nil, errors.New("no sites match when using id")
		}

		doc, err := bson.Marshal(&results[0])
		if err != nil {
			log.Errorf("1 %s", err)
			return nil, err
		}

		var agentCheck *Check
		err = bson.Unmarshal(doc, &agentCheck)
		if err != nil {
			log.Errorf("2 %s", err)
			return nil, err
		}

		ac.AgentID = agentCheck.AgentID
		ac.Type = agentCheck.Type
		ac.Duration = agentCheck.Duration
		ac.Server = agentCheck.Server
		ac.Triggered = agentCheck.Triggered
		ac.Count = agentCheck.Count
		ac.Pending = agentCheck.Pending
		ac.Target = agentCheck.Target

		return nil, nil
	} else {
		var agentChecks []*Check

		for _, r := range results {
			var acData Check
			doc, err := bson.Marshal(r)
			if err != nil {
				log.Errorf("1 %s", err)
				return nil, err
			}
			err = bson.Unmarshal(doc, &acData)
			if err != nil {
				log.Errorf("22 %s", err)
				return nil, err
			}

			agentChecks = append(agentChecks, &acData)
		}

		return agentChecks, nil
	}

	return nil, nil
}

// GetAll get all checks based on id, and &/or type
func (ac *Check) GetAll(db *mongo.Database) ([]*Check, error) {
	var filter = bson.D{{"agent", ac.AgentID}}
	if ac.Type != "" {
		filter = bson.D{{"agent", ac.AgentID}, {"type", ac.Type}}
	}

	cursor, err := db.Collection("agent_check").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	var agentCheck []*Check

	for _, rb := range results {
		m, err := bson.Marshal(&rb)
		if err != nil {
			log.Errorf("2 %s", err)
			return nil, err
		}
		var tC Check
		err = bson.Unmarshal(m, &tC)
		if err != nil {
			return nil, err
		}
		agentCheck = append(agentCheck, &tC)
	}
	return agentCheck, nil
}

func (ac *Check) Update(db *mongo.Database) error {
	var filter = bson.D{{"_id", ac.ID}}

	marshal, err := bson.Marshal(ac)
	if err != nil {
		return err
	}

	var b bson.D
	err = bson.Unmarshal(marshal, &b)
	if err != nil {
		log.Errorf("error unmarhsalling agent data when creating: %s", err)
		return err
	}

	update := bson.D{{"$set", b}}

	_, err = db.Collection("agent_check").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

//todo deleting checks