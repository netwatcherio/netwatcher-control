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
