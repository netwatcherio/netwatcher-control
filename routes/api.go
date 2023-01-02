package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/netwatcherio/netwatcher-agent/api"
	"github.com/netwatcherio/netwatcher-agent/checks"
	_ "github.com/netwatcherio/netwatcher-agent/checks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/handler"
)

func (r *Router) apiGetConfig() {
	r.App.Post("/api/v2/config/", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := api.Data{}

		var dataRequest api.Data

		err := json.Unmarshal(c.Body(), dataRequest)
		if err != nil {
			respB.Error = "500"
		}

		var agentSearch handler.Agent
		if dataRequest.ID != "" && dataRequest.PIN != "" {
			agentSearch.Pin = dataRequest.PIN
			hexId, err := primitive.ObjectIDFromHex(dataRequest.ID)
			if err != nil {
				respB.Error = "500"
			}
			agentSearch.ID = hexId
		} else if dataRequest.ID == "" && dataRequest.PIN != "" {
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
			respB.Checks = []checks.CheckData{}
			//todo add checks to be processed

			var agentCheck handler.AgentCheck

			all, err := agentCheck.GetAll(r.DB)
			if err != nil {
				return err
			}

			for _, ac := range all {
				modifiedData := checks.CheckData{
					Type:     string(ac.Type),
					Target:   ac.Target,
					ID:       ac.ID.Hex(),
					Duration: ac.Duration,
					Count:    ac.Count,
					Server:   ac.Server,
				}

				respB.Checks = append(respB.Checks, modifiedData)
			}
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			return c.SendString("Something went wrong...") // => ✋
		}
		return c.SendString(string(jRespB)) // => ✋ good
	})
}

func (r *Router) apiCheckData() {
	r.App.Post("/api/v2/check/", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := api.Data{}

		var dataRequest api.Data

		err := json.Unmarshal(c.Body(), dataRequest)
		if err != nil {
			respB.Error = "500"
		}

		if dataRequest.ID != "" && dataRequest.PIN != "" && len(dataRequest.Checks) > 0 {
			hexId, err := primitive.ObjectIDFromHex(dataRequest.ID)
			if err != nil {
				respB.Error = "500"
			}

			for _, cD := range dataRequest.Checks {
				checkId, err := primitive.ObjectIDFromHex(cD.ID)
				if err != nil {
					respB.Error = ""
				}

				data := handler.CheckData{
					Target:    cD.Target,
					ID:        primitive.NewObjectID(),
					CheckID:   checkId,
					AgentID:   hexId,
					Triggered: cD.Triggered,
					Result:    cD.Result,
				}
				data.Create(r.DB)
			}
		} else {
			respB.Error = "500"
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			return c.SendString("Something went wrong...") // => ✋
		}
		return c.SendString(string(jRespB)) // => ✋ good
	})
}
