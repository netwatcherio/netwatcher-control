package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sagostin/netwatcher-agent/agent_models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

	_, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
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

		var data []*agent_models.MtrTarget
		err = json.Unmarshal(dataJ, &data)
		if err != nil {
			log.Errorf("2 %s", err)
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

	_, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
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

		var data2 []*agent_models.NetworkInfo
		err = json.Unmarshal(dataJ, &data2)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}
		//log.Infof(string(dataJ))
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

	_, hash, b, err := verifyAgentHash(data.Pin, data.Hash, db)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		var data2 agent_models.SpeedTestInfo

		speedD, _ := json.Marshal(data.Data)
		if err != nil {
			log.Errorf("8 %s", err)
			respB.Response = 500
		}

		//log.Warnf("%s", c.Body())

		//fmt.Println(res["json"])

		//log.Infof("%s", string(jMar))
		err := json.Unmarshal(speedD, &data2)
		if err != nil {
			log.Errorf("2 %s", err)
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

		insertIcmpData(agent, data2, data.Timestamp, db)
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

	_, hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
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

		data := agent_models.AgentConfig{
			PingTargets:      nil,
			TraceTargets:     nil,
			PingInterval:     2,
			SpeedTestPending: true,
			TraceInterval:    5, // minutes
		}

		data.TraceTargets = append(data.TraceTargets, "1.1.1.1")
		data.PingTargets = append(data.PingTargets, "1.1.1.1")

		respB.Config = data
	}

	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}

	return ""
}
