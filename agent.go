package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sagostin/netwatcher-agent/agent_models"
	"github.com/sagostin/netwatcher-control/control_models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createAgent(c *mongo.Database) (bool, error) {
	var agentCfg = agent_models.AgentConfig{
		PingTargets:      []string{"1.1.1.1"},
		TraceTargets:     []string{"1.1.1.1"},
		PingInterval:     2,
		SpeedTestPending: true,
		TraceInterval:    5,
	}

	var agent = control_models.Agent{
		Site:        "dev",
		AgentConfig: agentCfg,
		Pin:         "12345",
		Hash:        "",
	}
	mar, err := bson.Marshal(agent)
	if err != nil {
		log.Errorf("1 %s", err)
		return false, err
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("2 %s", err)
		return false, err
	}
	result, err := c.Collection("agents").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return false, err
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return true, nil
}

// verifyAgent returns BOOL if verified, HASH (string to compare to inputted), BOOL if agent is new, error if something else
func verifyAgentHash(pin string, hash string, db *mongo.Database) (string, bool, error) {
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
		return "", false, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return "", false, err
	}

	fmt.Println(results)

	if len(results) > 1 {
		return "", false, errors.New("multiple agents match")
	}

	if len(results) == 0 {
		return "", false, errors.New("no agents match")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return "", false, err
	}

	var agent control_models.Agent
	err = bson.Unmarshal(doc, &agent)
	if err != nil {
		log.Errorf("2 %s", err)
		return "", false, err
	}

	if agent.Hash == "" {
		// todo generate hash and save to db
		return "test", true, nil
	} else if agent.Hash == hash && agent.Hash != "" {
		// return blank hash if verified
		return "", false, nil
	}

	return "", false, errors.New("something went wrong verifying agent")
}

func deleteAgent() {
	// todo remove from site
}
