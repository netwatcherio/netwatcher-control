package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/netwatcherio/netwatcher-agent/api"
	"github.com/netwatcherio/netwatcher-agent/checks"
	_ "github.com/netwatcherio/netwatcher-agent/checks"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/handler"
)

func (r *Router) apiGetConfig() {
	r.App.Post("/api/v2/config/", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := handler.Data{}

		var dataRequest handler.Data

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
			var cde []checks.CheckData

			for n := range all {
				ac := all[n]
				modifiedData := checks.CheckData{
					Type:     string(ac.Type),
					Target:   ac.Target,
					ID:       ac.ID,
					Duration: ac.Duration,
					Count:    ac.Count,
					Server:   ac.Server,
					Pending:  ac.Pending,
					Interval: ac.Interval,
				}

				cde = append(cde, modifiedData)
			}

			respB.Checks = cde
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
		respB := api.Data{}

		var dataRequest api.Data

		fmt.Println(string(c.Body()))

		err := json.Unmarshal(c.Body(), &dataRequest)
		if err != nil {
			respB.Error = "500 unable to read data"
			log.Fatal(err)
		}

		if dataRequest.ID != "000000000000000000000000" && dataRequest.PIN != "" && len(dataRequest.Checks) > 0 {
			hexId, err := primitive.ObjectIDFromHex(dataRequest.ID)
			if err != nil {
				respB.Error = "500 something went wrong, unable to compute object id"
			}

			for _, cD := range dataRequest.Checks {
				data := handler.CheckData{
					Target:    cD.Target,
					ID:        primitive.NewObjectID(),
					CheckID:   cD.ID,
					AgentID:   hexId,
					Triggered: cD.Triggered,
					Result:    cD.Result,
				}
				err = data.Create(r.DB)
				if err != nil {
					respB.Error = "500 unable to create check data"
				}

				if cD.Type == string(handler.CtSpeedtest) {
					ac := handler.AgentCheck{ID: cD.ID}
					_, err := ac.Get(r.DB)
					if err != nil {
						log.Error(err)
					}
					ac.Pending = false
					err = ac.Update(r.DB)
					if err != nil {
						log.Error(err)
					}
				}
			}
		} else {
			respB.Error = "500 unable to verify auth"
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			return c.SendString("Something went wrong...") // => ✋
		}
		return c.SendString(string(jRespB)) // => ✋ good
	})
}
