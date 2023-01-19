package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/netwatcherio/netwatcher-control/handler/agent"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
	"github.com/netwatcherio/netwatcher-control/handler/site"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) deleteAgent() {
	r.App.Get("/delete_agent/:agent?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("agent"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// delete agent
		a := agent.Agent{ID: aId}
		err = a.Delete(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// delete data
		ac := agent.Data{AgentID: aId}
		err = ac.Delete(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// delete check data
		acc := agent.Check{AgentID: aId}
		err = acc.Delete(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

func (r *Router) getAgent() {
	r.App.Get("/agent/:agent?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("agent"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		a := agent.Agent{ID: aId}
		err = a.Get(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(a)
	})
}

func (r *Router) getAgents() {
	r.App.Get("/agents/:siteid?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		a := site.Site{ID: aId}
		err = a.Get(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		agents, err := a.GetAgents(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(agents)
	})
}

func (r *Router) getGeneralAgentStats() {
	r.App.Get("/agent_stats/:agent?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("agent"))
		if err != nil {
			return c.JSON(err)
		}

		a := agent.Agent{ID: aId}
		err = a.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		stats, err := a.GetLatestStats(r.DB)
		if err != nil {
			log.Error(err)
		}

		return c.JSON(stats)
	})
}

func (r *Router) agentNew() {
	r.App.Post("/agent/new/:siteid?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		sId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		s := site.Site{ID: sId}
		err = s.Get(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		cAgent := new(agent.Agent)
		if err := c.BodyParser(cAgent); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		cAgent.Site = s.ID

		err = cAgent.Create(r.DB)
		if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}
