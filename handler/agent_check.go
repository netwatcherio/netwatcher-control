package handler

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

type AgentCheck struct {
	Type            CheckType          `json:"type"bson:"type""`
	Target          string             `json:"target,omitempty"bson:"target,omitempty"`
	ID              primitive.ObjectID `json:"id"bson:"_id"`
	AgentID         primitive.ObjectID `json:"agent"bson:"agent"`
	Duration        int                `json:"duration,omitempty'"bson:"duration,omitempty"`
	Count           int                `json:"count,omitempty"bson:"count,omitempty"`
	Triggered       bool               `json:"triggered"bson:"triggered,omitempty"`
	Server          bool               `json:"server,omitempty"bson:"server,omitempty"`
	Pending         bool               `json:"pending,omitempty"bson:"pending,omitempty"`
	CreateTimestamp time.Time          `bson:"create_timestamp"json:"create_timestamp"`
}

func (ac *AgentCheck) GetData(limit int64, justCheckId bool, recent bool, timeStart *time.Time, timeEnd *time.Time, db *mongo.Database) ([]*CheckData, error) {
	opts := options.Find().SetLimit(limit)
	var filter = bson.D{{"check", ac.ID}, {"type", ac.Type}}
	if ac.AgentID != (primitive.ObjectID{0}) {
		filter = bson.D{{"agent", ac.AgentID}, {"type", ac.Type}}
	}
	if justCheckId {
		filter = bson.D{{"check", ac.ID}}
	}

	var timeFilter bson.M

	if recent {
		opts = opts.SetSort(bson.D{{"timestamp", -1}})
	} else {
		timeFilter = bson.M{
			"check": ac.ID,
			"timestamp": bson.M{
				"$gt": timeStart,
				"$lt": timeEnd,
			}}
	}

	if timeStart != nil && timeEnd != nil {
		// todo change opts to use start and end time to check latest checks
	}

	var cursor *mongo.Cursor
	var err error

	if recent {
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

	var checkData []*CheckData

	for _, r := range results {
		var cData CheckData
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

func (ac *AgentCheck) Create(db *mongo.Database) error {
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

func (ac *AgentCheck) Get(db *mongo.Database) ([]*AgentCheck, error) {
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

		var agentCheck *AgentCheck
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
		var agentChecks []*AgentCheck

		for _, r := range results {
			var acData AgentCheck
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
func (ac *AgentCheck) GetAll(db *mongo.Database) ([]*AgentCheck, error) {
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
	var agentCheck []*AgentCheck

	for _, rb := range results {
		m, err := bson.Marshal(&rb)
		if err != nil {
			log.Errorf("2 %s", err)
			return nil, err
		}
		var tC AgentCheck
		err = bson.Unmarshal(m, &tC)
		if err != nil {
			return nil, err
		}
		agentCheck = append(agentCheck, &tC)
	}
	return agentCheck, nil
}

func (a *AgentCheck) Update(db *mongo.Database) error {
	var filter = bson.D{{"_id", a.ID}}

	marshal, err := bson.Marshal(&a)
	if err != nil {
		return err
	}

	var b *bson.D
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
