package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/netwatcherio/netwatcher-agent/checks"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Agent struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"json:"id"`
	Name        string             `bson:"name"json:"name"form:"name"`
	Site        primitive.ObjectID `bson:"site"json:"site"` // _id of mongo object
	Pin         string             `bson:"pin"json:"pin"`   // used for registration & authentication
	Heartbeat   time.Time          `bson:"heartbeat,omitempty"json:"heartbeat,omitempty"`
	Initialized bool               `bson:"initialized"json:"initialized"`
	Longitude   float64            `bson:"longitude"json:"longitude"form:"longitude"'`
	Latitude    float64            `bson:"latitude"json:"latitude"form:"latitude"`
	// pin can be regenerated, by setting hash blank, and when registering agents, it checks for blank hashs.
}

type AgentStats struct {
	AgentID          primitive.ObjectID `json:"agent_id"`
	Name             string             `json:"name"`
	Heartbeat        time.Time          `json:"heartbeat"`
	NetInfo          checks.NetResult   `json:"net_info"`
	SpeedTestInfo    checks.SpeedTest   `json:"speed_test_info"`
	SpeedTestPending bool               `json:"speed_test_pending"`
}

func (a *Agent) GetAgentStats(db *mongo.Database) (*AgentStats, error) {
	var stats AgentStats
	stats.AgentID = a.ID
	stats.Name = a.Name
	stats.Heartbeat = a.Heartbeat

	// get the latest net stats
	agentCheck := AgentCheck{AgentID: a.ID, Type: CT_NetInfo}
	get, err := agentCheck.GetData(1, true, nil, nil, db)
	if err != nil {
		return &stats, err
	}

	bytes, err := json.Marshal(get[0])
	if err != nil {
		return &stats, err
	}
	var netInfo checks.NetResult

	err = json.Unmarshal(bytes, &netInfo)
	if err != nil {
		return &stats, err
	}

	stats.NetInfo = netInfo

	// todo check the agent check itself to see if the speedtest is pending, else check and add the speedtest stats

	// get the latest net stats
	agentCheck = AgentCheck{AgentID: a.ID, Type: CT_SpeedTest}
	get, err = agentCheck.GetData(1, true, nil, nil, db)
	if err != nil {
		return &stats, err
	}

	bytes, err = json.Marshal(get[0])
	if err != nil {
		return &stats, err
	}
	var speedTest checks.SpeedTest

	err = json.Unmarshal(bytes, &speedTest)
	if err != nil {
		return &stats, err
	}
	stats.SpeedTestInfo = speedTest

	return &stats, nil
}

/*var agent = models.Agent{
	ID:   primitive.NewObjectID(),
	Site: site,
	Name: name,
	Initialized: false,
	Pin:  GeneratePin(9),
}*/

func (a *Agent) Create(db *mongo.Database) error {
	// todo handle to check if agent id is set and all that...
	a.Pin = GeneratePin(9)
	a.ID = primitive.NewObjectID()
	a.Initialized = false

	mar, err := bson.Marshal(a)
	if err != nil {
		log.Errorf("error marshalling agent data when creating: %s", err)
		return err
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("error unmarhsalling agent data when creating: %s", err)
		return err
	}
	result, err := db.Collection("agents").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("error inserting to database: %s", err)
		return err
	}

	fmt.Printf("created agent with id: %v\n", result.InsertedID)
	return nil
}

func (a *Agent) Get(db *mongo.Database) error {
	var filter = bson.D{{"_id", a.ID}}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		log.Errorf("error searching database for agent: %s", err)
		return err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Errorf("error cursoring through agents: %s", err)
		return err
	}

	//fmt.Println(results)

	if len(results) > 1 {
		return errors.New("multiple agents match when getting using id") // edge case??
	}

	if len(results) == 0 {
		return errors.New("no agents match when getting using id")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return err
	}

	var agent *Agent
	err = bson.Unmarshal(doc, &agent)
	if err != nil {
		log.Errorf("2 %s", err)
		return err
	}

	a.Latitude = agent.Latitude
	a.Longitude = agent.Longitude
	a.Pin = agent.Pin
	a.Name = agent.Name
	a.Site = agent.Site
	a.Heartbeat = agent.Heartbeat
	a.Initialized = agent.Initialized

	return nil
}

func (a *Agent) UpdateHeartBeat(db *mongo.Database) error {
	var filter = bson.D{{"_id", a.ID}}

	update := bson.D{{"$set", bson.D{{"heartbeat", time.Now()}}}}

	_, err := db.Collection("agents").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
