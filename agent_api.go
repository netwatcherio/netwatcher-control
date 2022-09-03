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

	hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
	if err != nil {
		log.Errorf("%s", err)
		respB.Response = 401
	}

	if hash != "" && b {

		var data []*agent_models.MtrTarget
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}
		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			return string(jRespB) // => ✋ good
		}
	}

	return ""
}

func apiUpdateNetwork(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiResponse{}
	respB.Response = 200

	var data *agent_models.Api
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		log.Errorf("2 %s", err)
		respB.Response = 500
	}

	hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
	if err != nil {
		log.Errorf("%s", err)
		respB.Response = 401
	}

	hash != "" && b

	jRespB, err := json.Marshal(respB)
	if err != nil {
		log.Errorf("3 Unable to marshal API response.")
	} else {
		return string(jRespB) // => ✋ good
	}
	return ""
}

func apiUpdateSpeedtest(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiResponse{}
	respB.Response = 200

	hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
	if err != nil {
		log.Errorf("%s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		var data agent_models.SpeedTestInfo

		//log.Warnf("%s", c.Body())

		//fmt.Println(res["json"])

		//log.Infof("%s", string(jMar))
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			return string(jRespB) // => ✋ good
		}
	}

	return ""
}

func apiUpdateIcmp(c *fiber.Ctx, db *mongo.Database) string {
	c.Accepts("Application/json") // "Application/json"
	respB := agent_models.ApiResponse{}
	respB.Response = 200

	hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
	if err != nil {
		log.Errorf("%s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		var data []*agent_models.IcmpTarget
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			return string(jRespB) // => ✋ good
		}
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

	hash, b, err := verifyAgentHash(c.Params("pin"), c.Params("hash"), db)
	if err != nil {
		log.Errorf("%s", err)
		respB.Response = 401
	}

	if hash != "" && b {
		respB.NewAgent = true
		respB.AgentHash = hash
		log.Infof("verified agent id %s", c.Params("pin"))
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
