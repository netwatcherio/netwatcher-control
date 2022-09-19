package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/netwatcherio/netwatcher-agent/agent_models"
	"github.com/netwatcherio/netwatcher-control/control_models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// CreateAgent returns error or object id of new agent
func CreateAgent(name string, icmpT []string, mtrT []string, site primitive.ObjectID, c *mongo.Database) (primitive.ObjectID, error) {
	var agentCfg = agent_models.AgentConfig{
		PingTargets:      icmpT,
		TraceTargets:     mtrT,
		PingInterval:     2,
		SpeedTestPending: true,
		TraceInterval:    5,
	}

	var agent = control_models.Agent{
		ID:          primitive.NewObjectID(),
		Site:        site,
		Name:        name,
		AgentConfig: agentCfg,
		Pin:         GeneratePin(9),
		Hash:        "",
	}
	mar, err := bson.Marshal(agent)
	if err != nil {
		log.Errorf("1 %s", err)
		return primitive.ObjectID{}, err
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("2 %s", err)
		return primitive.ObjectID{}, err
	}
	result, err := c.Collection("agents").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return primitive.ObjectID{}, err
	}

	fmt.Printf(" with _id: %v\n", result.InsertedID)
	return agent.ID, nil
}

func getAgents(site primitive.ObjectID, db *mongo.Database) ([]*control_models.Agent, error) {
	// if hash is blank, search for pin matching with blank hash
	// if none exist, return error
	// if match, return new agent, and new hash, then let another function update the hash?
	// if hash is included, search both, and return nil for hash, and false for new if verified
	// if hash is included and none match, return err
	var filter = bson.D{{"site", site}}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	//fmt.Println(results)

	if len(results) == 0 {
		return nil, errors.New("no agents match when using id")
	}

	var agent []*control_models.Agent
	for i := range results {
		doc, err := bson.Marshal(&results[i])
		if err != nil {
			log.Errorf("1 %s", err)
			return nil, err
		}
		var a *control_models.Agent
		err = bson.Unmarshal(doc, &a)
		if err != nil {
			log.Errorf("2 %s", err)
			return nil, err
		}

		agent = append(agent, a)
	}

	return agent, nil
}

func updateHeartbeat(id primitive.ObjectID, db *mongo.Database) error {
	var filter = bson.D{{"_id", id}}

	update := bson.D{{"$set", bson.D{{"heartbeat", time.Now()}}}}

	_, err := db.Collection("agents").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func getAgent(id primitive.ObjectID, db *mongo.Database) (*control_models.Agent, error) {
	// if hash is blank, search for pin matching with blank hash
	// if none exist, return error
	// if match, return new agent, and new hash, then let another function update the hash?
	// if hash is included, search both, and return nil for hash, and false for new if verified
	// if hash is included and none match, return err
	var filter = bson.D{{"_id", id}}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	//fmt.Println(results)

	if len(results) > 1 {
		return nil, errors.New("multiple agents match when using id")
	}

	if len(results) == 0 {
		return nil, errors.New("no agents match when using id")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return nil, err
	}

	var agent *control_models.Agent
	err = bson.Unmarshal(doc, &agent)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, err
	}

	err = updateHeartbeat(agent.ID, db)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

// verifyAgent returns string _id if has agent, HASH (string to compare to inputted), BOOL if agent is new, error if something else
func verifyAgentHash(pin string, hash string, db *mongo.Database) (primitive.ObjectID, string, bool, error) {
	// if hash is blank, search for pin matching with blank hash
	// if none exist, return error
	// if match, return new agent, and new hash, then let another function update the hash?
	// if hash is included, search both, and return nil for hash, and false for new if verified
	// if hash is included and none match, return err
	var filter = bson.D{{"pin", pin}, {"hash", ""}}
	if hash != "" {
		filter = bson.D{{"pin", pin}, {"hash", hash}}
	}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		return primitive.ObjectID{}, "", false, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return primitive.ObjectID{}, "", false, err
	}

	//fmt.Println(results)

	if len(results) > 1 {
		return primitive.ObjectID{}, "", false, errors.New("multiple agents match")
	}

	if len(results) == 0 {
		return primitive.ObjectID{}, "", false, errors.New("no agents match")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return primitive.ObjectID{}, "", false, err
	}

	var agent control_models.Agent
	err = bson.Unmarshal(doc, &agent)
	if err != nil {
		log.Errorf("2 %s", err)
		return primitive.ObjectID{}, "", false, err
	}

	if agent.Hash == "" {
		// todo generate hash and save to db
		return agent.ID, "test", true, nil
	} else if agent.Hash == hash && agent.Hash != "" {
		// return blank hash if verified
		return primitive.ObjectID{}, "", false, nil
	}

	return primitive.ObjectID{}, "", false, errors.New("something went wrong verifying agent")
}

func getLatestNetworkData(id primitive.ObjectID, db *mongo.Database) (control_models.NetworkData, error) {
	var filter = bson.D{{"agent", id}}
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
	cursor, err := db.Collection("network_data").Find(context.TODO(), filter, opts)
	if err != nil {
		return control_models.NetworkData{}, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return control_models.NetworkData{}, err
	}

	if len(results) == 0 {
		return control_models.NetworkData{}, errors.New("no agents match when using id")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return control_models.NetworkData{}, err
	}

	var netInfo1 control_models.NetworkData
	err = bson.Unmarshal(doc, &netInfo1)
	if err != nil {
		log.Errorf("22 %s", err)
		return control_models.NetworkData{}, err
	}

	return netInfo1, nil
}

func getIcmpData(id primitive.ObjectID, timeRange time.Duration, db *mongo.Database) ([]control_models.IcmpData, error) {
	var filter = bson.M{
		"agent": id,
		"timestamp": bson.M{
			"$gt": time.Now().Add(-timeRange),
			"$lt": time.Now(),
		}}

	cursor, err := db.Collection("icmp_data").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no data found")
	}

	var icmpD []control_models.IcmpData

	for _, r := range results {
		var icmp control_models.IcmpData
		doc, err := bson.Marshal(r)
		if err != nil {
			log.Errorf("1 %s", err)
			return nil, err
		}
		err = bson.Unmarshal(doc, &icmp)
		if err != nil {
			log.Errorf("22 %s", err)
			return nil, err
		}

		icmpD = append(icmpD, icmp)
	}

	/*j, err := json.Marshal(icmpD)
	if err != nil {
		log.Errorf("123 %s", err)
		return nil, err
	}
	log.Warnf("%s", j)*/

	return icmpD, nil
}

// getAgentStats get the general stats of agents from a site id objId
func getAgentStats(objId primitive.ObjectID, db *mongo.Database) ([]control_models.AgentStats, error) {
	site, err := getSite(objId, db)
	if err != nil {
		return nil, err
	}
	agents, err := getAgents(site.ID, db)
	if err != nil {
		return nil, err
	}
	var statsList []control_models.AgentStats

	for _, t := range agents {
		netInfo, err := getLatestNetworkData(t.ID, db)
		if err != nil {
			// todo handle error better
			//return nil, err
		}

		tNow := time.Now()
		tBeat := t.Heartbeat
		tDif := tNow.Sub(tBeat)

		var online = true
		if tDif.Minutes() > 2 {
			online = false
		}

		var agent = control_models.AgentStats{
			ID:          t.ID,
			Name:        t.Name,
			Heartbeat:   t.Heartbeat,
			NetworkInfo: netInfo.Data,
			LastSeen:    tDif,
			Online:      online,
		}
		statsList = append(statsList, agent)
	}
	return statsList, nil
}

func getAgentCount(site primitive.ObjectID, db *mongo.Database) (int, error) {
	var filter = bson.D{{"site", site}}

	count, err := db.Collection("agents").CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func deleteAgent() {
	// todo remove from site
}
