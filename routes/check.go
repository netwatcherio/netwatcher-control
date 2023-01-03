package routes

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/handler"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) check() {
	r.App.Get("/traceroute/:mtrid?", func(c *fiber.Ctx) error {
		user, err := validateUser(r, c)
		if err != nil {
			return c.Redirect("/auth")
		}

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
		/*var agentStatList models.AgentStatsList

		  stats, err := getAgentStatsForSite(objId, db)
		  if err != nil {
		  	//todo handle error
		  	//return err
		  }
		  agentStatList.List = stats

		  var hasAgents = true
		  if len(agentStatList.List) == 0 {
		  	hasAgents = false
		  }

		  doc, err := json.Marshal(agentStatList)
		  if err != nil {
		  	log.Errorf("1 %s", err)
		  }*/

		/*agents, err := getAgents(objId, db)
		  if err != nil {
		  	// todo handle error
		  	//return err
		  }

		  doc, err := json.Marshal(agents)
		  if err != nil {
		  	log.Errorf("1 %s", err)
		  }

		  var hasAgentsBool = true
		  if len(agents) == 0 {
		  	hasAgentsBool = false
		  	log.Warnf("%s", "site does NOT have agents")
		  }*/

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
			"email":        user.Email},
			/*"agents":       html.UnescapeString(string(doc)),
			  "hasAgents":    hasAgents},*/
			"layouts/main")
	})
}

func (r *Router) checkNew() {
	r.App.Get("/check/new/:agentid?", func(c *fiber.Ctx) error {
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
			return c.Redirect("/home")
		}

		site := handler.Site{ID: agent.Site}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("check_new", fiber.Map{
			"title":        "new check",
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

		check := handler.AgentCheck{AgentID: cAgent.ID}
		check.Type = handler.CT_NetInfo
		err = check.Create(r.DB)
		if err != nil {
			return err
		}

		// todo create default checks such as network info and that sort of thing

		// todo handle error/success and return to home also display message for error if error
		return c.Redirect("/agent/install/" + cAgent.ID.Hex())
	})
}
