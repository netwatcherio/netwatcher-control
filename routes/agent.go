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

func (r *Router) getAgent() {
	r.App.Get("/agent/:agent?", func(c *fiber.Ctx) error {
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

		return c.JSON(a)
	})
}

func (r *Router) getGeneralAgentStats() {
	r.App.Get("/agent/:agent?", func(c *fiber.Ctx) error {
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
			return c.JSON(err)
		}

		sId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.JSON(err)
		}

		s := site.Site{ID: sId}
		err = s.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		cAgent := new(agent.Agent)
		if err := c.BodyParser(cAgent); err != nil {
			return c.JSON(err)
		}

		cAgent.Site = s.ID

		err = cAgent.Create(r.DB)
		if err != nil {
			log.Error(err)
			return c.JSON(err)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}
