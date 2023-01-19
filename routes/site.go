package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/netwatcherio/netwatcher-control/handler/agent"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
	"github.com/netwatcherio/netwatcher-control/handler/site"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) addSiteMember() {
	r.App.Post("/site/:siteid?/members/add", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		m := new(site.AddSiteMember)
		if err := c.BodyParser(m); err != nil {
			return c.JSON(err)
		}

		mU := auth.User{Email: m.Email}
		email, err := mU.FromEmail(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		siteId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.JSON(err)
		}

		s := site.Site{ID: siteId}
		err = s.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		err = s.AddMember(email.ID, m.Role, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		err = email.AddSite(s.ID, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

func (r *Router) newSite() {
	r.App.Post("/site/new", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		u, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		s := new(site.Site)
		if err := c.BodyParser(s); err != nil {
			return c.JSON(err)
		}

		err = s.Create(u.ID, r.DB)
		if err != nil {
			return c.JSON(err)
		}
		return c.SendStatus(fiber.StatusOK)
	})
}
func (r *Router) getSite() {
	r.App.Get("/site/:siteid?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		siteId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.JSON(err)
		}

		s := site.Site{ID: siteId}
		err = s.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(s)
	})
}

func (r *Router) deleteSite() {
	r.App.Get("/delete_site/:site?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("site"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// delete site
		s := site.Site{ID: aId}
		agents, err := s.GetAgents(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		err = s.Delete(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, aa := range agents {
			a := agent.Agent{ID: aa.ID}
			err = a.Delete(r.DB)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			// delete data
			ac := agent.Data{AgentID: aa.ID}
			err = ac.Delete(r.DB)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			// delete check data
			acc := agent.Check{AgentID: aa.ID}
			err = acc.Delete(r.DB)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

// get braden
func (r *Router) getSites() {
	r.App.Get("/sites", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		user, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		type AgentCountInfo struct {
			SiteID primitive.ObjectID `json:"site_id"`
			Count  int                `json:"count"`
		}

		var sitesList struct {
			Sites          []*site.Site     `json:"sites"`
			AgentCountInfo []AgentCountInfo `json:"agent_counts"`
		}

		for _, sid := range user.Sites {
			s := site.Site{ID: sid}
			err := s.Get(r.DB)
			if err != nil {
				return c.JSON(err)
			}

			count, err := s.AgentCount(r.DB)
			if err != nil {
				return c.JSON(err)
			}

			tempCount := AgentCountInfo{
				SiteID: s.ID,
				Count:  count,
			}

			sitesList.Sites = append(sitesList.Sites, &s)
			sitesList.AgentCountInfo = append(sitesList.AgentCountInfo, tempCount)
		}

		return c.JSON(sitesList)
	})
}
