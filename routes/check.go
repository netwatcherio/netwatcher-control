package routes

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/handler"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) check() {
	r.App.Get("/check/:checkid?", func(c *fiber.Ctx) error {
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
			"agentId":      agent.ID.Hex(),
			"siteName":     site.Name,
			"siteSelected": true,
		},
			"layouts/main")
	})
	r.App.Post("/check/new/:agentid?", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		_, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("agentid") == "" {
			return c.Redirect("/agents")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("agentid"))
		if err != nil {
			return c.Redirect("/agents")
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

		cCheck := new(handler.CheckNewForm)
		if err := c.BodyParser(cCheck); err != nil {
			log.Warnf("4 %s", err)
			return err
		}

		var aC handler.AgentCheck

		// todo validate target depending on check
		// ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$

		if cCheck.Type == string(handler.CtMtr) {
			if cCheck.Duration < 30 {
				cCheck.Duration = 30
			}

			aC = handler.AgentCheck{
				Type:     handler.CtMtr,
				Target:   cCheck.Target,
				AgentID:  agent.ID,
				Duration: cCheck.Duration,
			}
		} else if cCheck.Type == string(handler.CtRperf) {
			aC = handler.AgentCheck{
				Type:    handler.CtRperf,
				Target:  cCheck.Target,
				AgentID: agent.ID,
				Server:  cCheck.RperfServerEnable,
			}
		}

		err = aC.Create(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/agents")
		}

		// todo create default checks such as network info and that sort of thing

		// todo handle error/success and return to home also display message for error if error
		return c.Redirect("/agent/" + agent.ID.Hex())
	})
}
