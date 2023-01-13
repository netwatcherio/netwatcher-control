package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/netwatcherio/netwatcher-control/handler"
	"github.com/netwatcherio/netwatcher-control/handler/agent/checks"
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
	Heartbeat   time.Time          `bson:"heartbeat,omitempty"json:"heartbeat"`
	Initialized bool               `bson:"initialized"json:"initialized"`
	Longitude   float64            `bson:"longitude"json:"longitude"form:"longitude"'`
	Latitude    float64            `bson:"latitude"json:"latitude"form:"latitude"`
	Timestamp   time.Time          `bson:"timestamp"json:"timestamp"`
	// pin can be regenerated, by setting hash blank, and when registering agents, it checks for blank hashs.
}

type Stats struct {
	// not used in api / db
	AgentID          primitive.ObjectID `json:"agent_id"`
	Name             string             `json:"name"`
	Heartbeat        time.Time          `json:"heartbeat"`
	NetInfo          checks.NetResult   `json:"net_info"`
	SpeedTestInfo    checks.SpeedTest   `json:"speed_test_info"`
	SpeedTestPending bool               `json:"speed_test_pending"`
}

func (a *Agent) GetLatestStats(db *mongo.Database) (*Stats, error) {
	var stats Stats
	stats.AgentID = a.ID
	stats.Name = a.Name
	stats.Heartbeat = a.Heartbeat

	// get the latest net stats
	agentCheck := checks.AgentCheck{AgentID: a.ID, Type: checks.CtNetinfo}
	netInfo, err := agentCheck.GetData(1, false, true, time.Time{}, time.Time{}, db)
	if err != nil {
		return &stats, err
	}

	netB, err := bson.Marshal(netInfo[0].Result)
	if err != nil {
		return &stats, err
	}

	var nf checks.NetResult
	err = bson.Unmarshal(netB, &nf)
	if err != nil {
		return &stats, err
	}

	stats.NetInfo.DefaultGateway = nf.DefaultGateway
	stats.NetInfo.InternetProvider = nf.InternetProvider
	stats.NetInfo.Timestamp = nf.Timestamp
	stats.NetInfo.LocalAddress = nf.LocalAddress
	stats.NetInfo.PublicAddress = nf.PublicAddress
	stats.NetInfo.Long = nf.Long
	stats.NetInfo.Lat = nf.Lat

	// todo check the agent check itself to see if the speedtest is pending, else check and add the speedtest stats

	// get the latest net stats
	agentCheck = checks.AgentCheck{AgentID: a.ID, Type: checks.CtSpeedtest}
	speedGet, err := agentCheck.GetData(1, false, true, time.Time{}, time.Time{}, db)
	if err != nil {
		return &stats, err
	}

	speedB, err := bson.Marshal(speedGet[0].Result)
	if err != nil {
		return &stats, err
	}

	var speedTest checks.SpeedTest

	err = bson.Unmarshal(speedB, &speedTest)
	if err != nil {
		return &stats, err
	}

	stats.SpeedTestInfo.DLSpeed = speedTest.DLSpeed
	stats.SpeedTestInfo.ULSpeed = speedTest.ULSpeed
	stats.SpeedTestInfo.Timestamp = speedTest.Timestamp
	stats.SpeedTestInfo.Server = speedTest.Server
	stats.SpeedTestInfo.Host = speedTest.Host
	stats.SpeedTestInfo.Latency = speedTest.Latency

	return &stats, err
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
	a.Pin = handler.GeneratePin(9)
	a.ID = primitive.NewObjectID()
	a.Initialized = false
	a.Timestamp = time.Now()

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

func (a *Agent) Verify(db *mongo.Database) error {
	var filter = bson.D{{"pin", a.Pin}, {"initialized", false}}
	initialized := false
	if &a.ID == nil {
		filter = bson.D{{"pin", a.Pin}, {"_id", a.ID}}
		initialized = true
	}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return err
	}

	if len(results) > 1 {
		return errors.New("multiple agents match")
	}

	if len(results) <= 0 {
		return errors.New("no agents match")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return err
	}

	var agent Agent
	err = bson.Unmarshal(doc, &agent)
	if err != nil {
		log.Errorf("2 %s", err)
		return err
	}
	a.ID = agent.ID
	a.Latitude = agent.Latitude
	a.Longitude = agent.Longitude
	a.Pin = agent.Pin
	a.Name = agent.Name
	a.Site = agent.Site
	a.Heartbeat = agent.Heartbeat

	if !initialized {
		a.Initialized = true
	} else {
		a.Initialized = agent.Initialized
	}

	err = a.Update(db)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) UpdateHeartbeat(db *mongo.Database) error {
	var filter = bson.D{{"_id", a.ID}}

	update := bson.D{{"$set", bson.D{{"heartbeat", time.Now()}}}}

	_, err := db.Collection("agents").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) Update(db *mongo.Database) error {
	var filter = bson.D{{"_id", a.ID}}

	marshal, err := bson.Marshal(a)
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

	_, err = db.Collection("agents").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
