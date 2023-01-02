package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html"
	"math"
	"netwatcher-control/handler"
	_ "strings"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) agent() {
	r.App.Get("/agent/:agent?", func(c *fiber.Ctx) error {
		b, _ := handler.ValidateSession(c, r.Session, r.DB)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := handler.GetUserFromSession(c, r.Session, r.DB)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		if c.Params("agent") == "" {
			return c.RedirectBack("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("agent"))
		if err != nil {
			return c.RedirectBack("/home")
		}

		agent := handler.Agent{ID: objId}
		err = agent.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/agents")
		}

		site := handler.Site{ID: objId}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		marshal, err := json.Marshal(agent)
		if err != nil {
			log.Errorf("13 %s", err)
		}

		getAgentStats, err := agent.GetAgentStats(r.DB)
		if err != nil {
			return err
		}

		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("agent", fiber.Map{
			"title":        agent.Name,
			"siteSelected": true,
			"siteName":     site.Name,
			"siteId":       site.ID.Hex(),
			"agents":       html.UnescapeString(string(marshal)),
			/*"mtr":              html.UnescapeString(string(marshalMtr)),
			"speed":            html.UnescapeString(string(marshalSpeed)),*/
			/*"online":           stats.Online,*/
			/*"agentStats":       html.UnescapeString(string(getAgentStats)),*/
			"publicAddress":    getAgentStats.NetInfo.PublicAddress,
			"localAddress":     getAgentStats.NetInfo.LocalAddress,
			"defaultGateway":   getAgentStats.NetInfo.DefaultGateway,
			"internetProvider": getAgentStats.NetInfo.InternetProvider,
			"uploadSpeed":      math.Round(getAgentStats.SpeedTestInfo.ULSpeed),
			"downloadSpeed":    math.Round(getAgentStats.SpeedTestInfo.DLSpeed),
			/*"speedtestPending": agent.AgentConfig.SpeedTestPending,*/
			"agentId":   agent.ID.Hex(),
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			//"icmpMetrics":  html.UnescapeString(string(j)),
		},
			"layouts/main")
	})
}

func (r *Router) agents() {
	r.App.Get("/agents/:siteid?", func(c *fiber.Ctx) error {
		user, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("siteid") == "" {
			return c.Redirect("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.Redirect("/home")
		}

		site := handler.Site{ID: objId}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		agentStats, err := site.GetAgentSiteStats(r.DB)
		if err != nil {
			log.Error(err)
		}

		marshalAS, err := json.Marshal(agentStats)
		if err != nil {
			return err
		}

		hasAgents := false
		if len(agentStats) > 0 {
			hasAgents = true
		}

		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		//log.Errorf("%s", string(doc))
		return c.Render("agents", fiber.Map{
			"title":        "agents",
			"siteSelected": true,
			"siteId":       site.ID.Hex(),
			"siteName":     site.Name,
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"agents":       html.UnescapeString(string(marshalAS)),
			"hasAgents":    hasAgents},
			"layouts/main")
	})
}

func (r *Router) agentNew() {
	r.App.Get("/agent/new/:siteid?", func(c *fiber.Ctx) error {
		user, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("siteid") == "" {
			return c.Redirect("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.Redirect("/home")
		}

		site := handler.Site{ID: objId}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("agent_new", fiber.Map{
			"title":        "new agent",
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"siteId":       site.ID.Hex(),
			"siteName":     site.Name,
			"siteSelected": true,
		},
			"layouts/main")
	})
	r.App.Post("/agent/new/:siteid?", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		_, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("siteid") == "" {
			return c.Redirect("/agents")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.Redirect("/agents")
		}

		site := handler.Site{ID: objId}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		cAgent := new(handler.Agent)
		if err := c.BodyParser(cAgent); err != nil {
			log.Warnf("4 %s", err)
			return err
		}

		cAgent.Site = site.ID

		err = cAgent.Create(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/agents")
		}

		// todo create default checks such as network info and that sort of thing

		// todo handle error/success and return to home also display message for error if error
		return c.Redirect("/agent/install/" + cAgent.ID.Hex())
	})
}

func (r *Router) agentInstall() {
	r.App.Get("/agent/install/:agentid", func(c *fiber.Ctx) error {
		user, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("agentid") == "" {
			return c.Redirect("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("agentid"))
		if err != nil {
			return c.Redirect("/home")
		}

		agent := handler.Agent{ID: objId}
		err = agent.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/agents")
		}

		site := handler.Site{ID: agent.Site}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		//todo handle if already installed

		return c.Render("agent_install", fiber.Map{
			"title":        "agent install",
			"siteSelected": true,
			"siteId":       agent.Site.Hex(),
			"siteName":     site.Name,
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"agentPin":     agent.Pin,
		},
			"layouts/main")
	})
}
