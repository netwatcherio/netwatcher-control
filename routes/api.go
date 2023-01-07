package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/netwatcherio/netwatcher-control/handler"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Router) apiGetConfig() {
	r.App.Post("/api/v2/config/", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := handler.ApiRequest{}

		var dataRequest handler.ApiRequest

		err := json.Unmarshal(c.Body(), &dataRequest)
		if err != nil {
			respB.Error = "500"
		}

		var agentSearch handler.Agent
		if dataRequest.ID != "000000000000000000000000" && dataRequest.PIN != "" {
			agentSearch.Pin = dataRequest.PIN
			hexId, err := primitive.ObjectIDFromHex(dataRequest.ID)
			if err != nil {
				respB.Error = "500"
			}
			agentSearch.ID = hexId
		} else if dataRequest.ID == "000000000000000000000000" && dataRequest.PIN != "" {
			agentSearch.Pin = dataRequest.PIN
			agentSearch.Initialized = false
		} else {
			respB.Error = "500"
		}

		if (agentSearch.ID != primitive.ObjectID{0} || agentSearch.Initialized == false) && (agentSearch.Pin != "") {
			err := agentSearch.Verify(r.DB)
			if err != nil {
				respB.Error = "500"
			}

			respB.ID = agentSearch.ID.Hex()
			respB.PIN = agentSearch.Pin
			//todo add checks to be processed

			agentCheck := handler.AgentCheck{
				AgentID: agentSearch.ID,
			}

			all, err := agentCheck.GetAll(r.DB)
			if err != nil {
				respB.Error = "500"
			}

			respB.Data = all
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			return c.SendString("Something went wrong...") // => ✋
		}
		return c.SendString(string(jRespB)) // => ✋ good
	})
}

func (r *Router) apiDataPush() {
	r.App.Post("/api/v2/agent/push", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := handler.ApiRequest{}

		var dataRequest handler.ApiRequest

		err := json.Unmarshal(c.Body(), &dataRequest)
		if err != nil {
			respB.Error = "500"
		}

		var agentSearch handler.Agent
		if dataRequest.ID != "000000000000000000000000" && dataRequest.PIN != "" {
			agentSearch.Pin = dataRequest.PIN
			hexId, err := primitive.ObjectIDFromHex(dataRequest.ID)
			if err != nil {
				respB.Error = "500"
			}
			agentSearch.ID = hexId
		} else if dataRequest.ID == "000000000000000000000000" && dataRequest.PIN != "" {
			agentSearch.Pin = dataRequest.PIN
			agentSearch.Initialized = false
		} else {
			respB.Error = "500"
		}

		if (agentSearch.ID != (primitive.ObjectID{0})) && (agentSearch.Pin != "") {
			err := agentSearch.Verify(r.DB)
			if err != nil {
				respB.Error = "500 auth failed"
			}

			respB.ID = agentSearch.ID.Hex()
			respB.PIN = agentSearch.Pin
			//todo add checks to be processed

			var checkD []handler.CheckData
			err = json.Unmarshal([]byte(dataRequest.Data.(string)), &checkD)
			if err != nil {
				log.Error(err)
			}

			respB.Error = ""

			for _, cd := range checkD {
				err := cd.Create(r.DB)
				if err != nil {
					log.Error(err)
				}
			}
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			return c.SendString("Something went wrong...") // => ✋
		}
		return c.SendString(string(jRespB)) // => ✋ good
	})
}
