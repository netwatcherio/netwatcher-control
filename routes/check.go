package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/handler"
	"netwatcher-control/models"
)

app.Get("/icmp/:agent?", func(c *fiber.Ctx) error {
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
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

	agent := models.Agent{ID: objId}
	err = agent.Get(db)
	if err != nil {
		log.Error(err)
		return c.Redirect("/agents")
	}

	site := models.Site{ID: objId}
	err = site.Get(db)
	if err != nil {
		log.Error(err)
		return c.Redirect("/home")
	}

	doc, err := json.Marshal(agent)
	if err != nil {
		log.Errorf("13 %s", err)
		return err
	}

	/*icmpD, err := getIcmpData(objId, time.Hour*24, db)
	if err != nil {
		return err
	}*/

	/*for n := range icmpD {
		for n2 := range icmpD[n].Data {
			icmpD[n].Data[n2].Result = nil
		}
	}*/

	/*j, err := json.Marshal(icmpD)
	if err != nil {
		log.Errorf("%s", err)
		return err
	}*/

	// Render index within layouts/main
	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("icmp", fiber.Map{
		"title":        agent.Name,
		"siteSelected": true,
		"siteName":     site.Name,
		"siteId":       site.ID.Hex(),
		"agents":       html.UnescapeString(string(doc)),
		/*"icmpMetrics":  html.UnescapeString(string(j)),*/
	},
		"layouts/main")
})

app.Get("/traceroute/:mtrid?", func(c *fiber.Ctx) error {
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
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

	agent := models.Agent{ID: objId}
	err = agent.Get(db)
	if err != nil {
		log.Error(err)
		return c.Redirect("/agents")
	}

	site := models.Site{ID: objId}
	err = site.Get(db)
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
app.Get("/traceroutes/:agent?", func(c *fiber.Ctx) error {
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
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

	agent := models.Agent{ID: objId}
	err = agent.Get(db)
	if err != nil {
		log.Error(err)
		return c.Redirect("/agents")
	}

	site := models.Site{ID: objId}
	err = site.Get(db)
	if err != nil {
		log.Error(err)
		return c.Redirect("/home")
	}

	/*mtrData, err := getLatestMtrData(agent.ID, 10, db)
	if err != nil {
		log.Errorf("15 %s", err)
	}*/

	/*marshalMtr, err := json.Marshal(mtrData)
	if err != nil {
		log.Errorf("16 %s", err)
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
	return c.Render("traceroutes", fiber.Map{
		"title":        agent.Name,
		"siteSelected": true,
		"siteId":       site.ID.Hex(),
		"siteName":     site.Name,
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email,
		"agentName":    agent.Name,
		/*"mtr":          html.UnescapeString(string(marshalMtr)),*/
	},
		"layouts/main")
})

app.Get("/alerts/:siteid?", func(c *fiber.Ctx) error {
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	if c.Params("siteid") == "" {
		return c.Redirect("/home")
	}
	objId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
	if err != nil {
		return c.Redirect("/home")
	}

	site := models.Site{ID: objId}
	err = site.Get(db)
	if err != nil {
		log.Error(err)
		return c.Redirect("/home")
	}

	// Render index within layouts/main
	// TODO process if they are logged in or not, otherwise send them to registration/login
	//log.Errorf("%s", string(doc))
	return c.Render("alerts", fiber.Map{
		"title":        "alerts",
		"siteSelected": true,
		"siteId":       site.ID.Hex(),
		"siteName":     site.Name,
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email},
		"layouts/main")
})
