package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"netwatcher-control/agent_models"
)

func apiUpdateMtr(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiConfigResponse{}
	respB.Response = 200

	var data *agent_models.ApiPushData
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 500
	}
	log.Errorf("%s", c.Body())

	agentId, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
	if err != nil {
		log.Errorf("0 %s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		dataJ, err := json.Marshal(data.Data)
		if err != nil {
			log.Errorf("%s 1", err)
			respB.Response = 401
		}

		//log.Infof("4 %s", dataJ)

		var data2 []agent_models.MtrTarget
		err = json.Unmarshal(dataJ, &data2)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		var agent, _ = getAgent(agentId, db)
		if err != nil {
			log.Errorf("5 %s", err)
			respB.Response = 500
		}

		_, err = insertMtrData(agent, data2, data.Timestamp, db)
		if err != nil {
			respB.Response = 500
		}

		err = updateHeartbeat(agent.ID, db)
		if err != nil {
			respB.Response = 500
		}
	}
	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}
	return ""
}

func apiUpdateNetwork(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiResponse{}
	respB.Response = 200

	var data *agent_models.ApiPushData
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 500
	}

	agentId, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
	if err != nil {
		log.Errorf("%s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		dataJ, err := json.Marshal(data.Data)
		if err != nil {
			log.Errorf("%s", err)
			respB.Response = 401
		}

		var data2 agent_models.NetworkInfo
		err = json.Unmarshal(dataJ, &data2)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		var agent, _ = getAgent(agentId, db)
		if err != nil {
			log.Errorf("5 %s", err)
			respB.Response = 500
		}

		_, err = insertNetworkInfo(agent, data2, data.Timestamp, db)
		if err != nil {
			respB.Response = 500
		}

		err = updateHeartbeat(agent.ID, db)
		if err != nil {
			respB.Response = 500
		}
	}

	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}
	return ""
}

func apiUpdateSpeedTest(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiResponse{}
	respB.Response = 200

	var data *agent_models.ApiPushData
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 500
	}

	agentId, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		speedD, _ := json.Marshal(data.Data)
		if err != nil {
			log.Errorf("8 %s", err)
			respB.Response = 500
		}

		var data2 agent_models.SpeedTestInfo
		err := json.Unmarshal(speedD, &data2)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		var agent, _ = getAgent(agentId, db)
		if err != nil {
			log.Errorf("5 %s", err)
			respB.Response = 500
		}

		_, err = insertSpeedTestData(agent, data2, data.Timestamp, db)
		if err != nil {
			respB.Response = 500
		}
		err = updateHeartbeat(agent.ID, db)
		if err != nil {
			respB.Response = 500
		}
	}

	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}

	return ""
}

func apiUpdateIcmp(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiResponse{}
	respB.Response = 200

	var data agent_models.ApiPushData
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 500
	}

	agentId, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
	if err != nil {
		log.Errorf("0 %s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		icmpD, _ := json.Marshal(data.Data)
		if err != nil {
			log.Errorf("1 %s", err)
			respB.Response = 500
		}

		var data2 []agent_models.IcmpTarget
		err := json.Unmarshal(icmpD, &data2)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		var agent, _ = getAgent(agentId, db)
		if err != nil {
			log.Errorf("5 %s", err)
			respB.Response = 500
		}

		_, err = insertIcmpData(agent, data2, data.Timestamp, db)
		if err != nil {
			respB.Response = 500
		}
		err = updateHeartbeat(agent.ID, db)
		if err != nil {
			respB.Response = 500
		}
	}

	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}

	return ""
}

func apiGetConfig(c *fiber.Ctx, db *mongo.Database) string {
	// todo change api get to post to provide auth inside of the post request, and return config as response

	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiConfigResponse{}
	respB.Response = 200

	if c.Params("pin") == "" {
		respB.Response = 401
	}

	agentId, hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
	if err != nil {
		log.Errorf("0 %s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		respB.NewAgent = true
		respB.AgentHash = hash
		//log.Infof("verified agent id %s", c.Params("pin"))
	}

	// verify agent and generate hash on configuration request
	// pin will be 9 characters and verified with hash
	// if no hash is included in the api path, then it will check if there are any
	// agents with a blank hash

	if respB.Response == 200 {

		var agent, _ = getAgent(agentId, db)
		if err != nil {
			log.Errorf("5 %s", err)
			respB.Response = 500
		}

		respB.Config = agent.AgentConfig
	}

	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}

	return ""
}
