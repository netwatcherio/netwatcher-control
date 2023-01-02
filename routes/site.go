package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"netwatcher-control/handler"
	"netwatcher-control/models"
)

app.Get("/site/:siteid?/members/add", func(c *fiber.Ctx) error {
	// Render index within layouts/main
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

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

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	//todo get agent count

	// convert to json for testing
	//siteJs, err := json.Marshal(sitesList)
	if err != nil {
		// todo handle properly
		return c.Redirect("/auth")
	}

	//log.Infof("%s", siteJs)

	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("site_member_add", fiber.Map{
		"title":        "add member",
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email,
		"siteSelected": true,
		"siteName":     site.Name,
		"siteId":       site.ID.Hex(),
		//"sites":     html.UnescapeString(string(siteJs)),
	},
		"layouts/main")
})
app.Post("/site/:siteid?/members/add", func(c *fiber.Ctx) error {
	c.Accepts("application/x-www-form-urlencoded") // "Application/json"
	// Render index within layouts/main
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

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

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	newMember := new(models.AddSiteMember)
	if err := c.BodyParser(newMember); err != nil {
		//todo
		//return err
	}

	if newMember.Role > 2 {
		// check if data has been tampered to make a new member the owner
		log.Warnf(" %s", "someone is trying to tamper with the roles when adding")
		// TODO support owner transferring
		return c.Redirect("/site/" + site.ID.Hex() + "/members")
	}

	var usrTmp = models.User{Email: newMember.Email}

	usr, err := usrTmp.GetUserFromEmail(db)
	if err != nil {
		log.Errorf("12 %s", err)
		//TODO handle error correctly
		return c.Redirect("/site/" + site.ID.Hex() + "")
	}

	b, err = site.AddMember(usr.ID, newMember.Role, db)
	if err != nil {
		log.Errorf("2 %s", err)
		//todo handle better
		return c.Redirect("/site/" + site.ID.Hex() + "")
	}

	if !b {
		log.Errorf("something went wrong adding member to site")
		return c.Redirect("/site/" + site.ID.Hex() + "")
	}
	addSite, err := usr.AddSite(site.ID, db)
	if err != nil {
		return err
	}
	if !addSite {
		log.Infof("%s", "somethiung went wrongies")
		return c.Redirect("/site/" + site.ID.Hex() + "")
	}
	log.Infof("%s", "added member to site successfully")

	// todo handle error/success and return to home
	return c.Redirect("/site/" + site.ID.Hex() + "/members")
})
app.Get("/site/:siteid?/members", func(c *fiber.Ctx) error {
	// Render index within layouts/main
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

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

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	var siteMembers []models.SiteMember
	for _, mem := range site.Members {
		siteMembers = append(siteMembers, mem)
	}

	var siteUsers []*models.User
	for _, usr := range siteMembers {
		c2 := &models.User{ID: usr.User}
		u, err := c2.GetUserFromID(db)
		if err != nil {
			log.Errorf("%s %s", "0 Error processing users in site id", site.ID.Hex())
		}
		siteUsers = append(siteUsers, u)
	}

	siteMem, err := json.Marshal(siteMembers)
	if err != nil {
		log.Errorf("%s %s", " Error processing members in site id", site.ID.Hex())
	}
	siteUsr, err := json.Marshal(siteUsers)
	if err != nil {
		log.Errorf("%s %s", "2 Error processing users in site id", site.ID.Hex())
	}

	//todo get agent count

	// convert to json for testing
	//siteJs, err := json.Marshal(sitesList)
	if err != nil {
		// todo handle properly
		return c.Redirect("/auth")
	}

	//log.Infof("%s", siteJs)

	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("site_members", fiber.Map{
		"title":        "members",
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email,
		"siteSelected": true,
		"siteName":     site.Name,
		"siteId":       site.ID.Hex(),
		"siteMem":      html.UnescapeString(string(siteMem)),
		"siteUsr":      html.UnescapeString(string(siteUsr)),
		//"sites":     html.UnescapeString(string(siteJs)),
	},
		"layouts/main")
})

// sites
app.Get("/sites", func(c *fiber.Ctx) error {
	// Render index within layouts/main
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	type AgentCountInfo struct {
		SiteID primitive.ObjectID `json:"site_id"`
		Count  int                `json:"count"`
	}

	var sitesList struct {
		Sites          []*models.Site   `json:"sites"`
		AgentCountInfo []AgentCountInfo `json:"agentCountInfo"`
	}

	for _, sid := range user.Sites {
		site := models.Site{ID: sid}
		err = site.Get(db)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		count, err := site.AgentCount(db)
		if err != nil {
			//todo handle error
		}

		tempCount := AgentCountInfo{
			SiteID: site.ID,
			Count:  count,
		}

		sitesList.Sites = append(sitesList.Sites, &site)
		sitesList.AgentCountInfo = append(sitesList.AgentCountInfo, tempCount)
	}

	var hasSites = true
	if len(user.Sites) == 0 {
		hasSites = false
	}

	//todo get agent count

	// convert to json for testing
	siteJs, err := json.Marshal(sitesList)
	if err != nil {
		// todo handle properly
		return c.Redirect("/auth")
	}

	log.Infof("%s", siteJs)

	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("sites", fiber.Map{
		"title":     "sites",
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"hasSites":  hasSites,
		"sites":     html.UnescapeString(string(siteJs)),
	},
		"layouts/main")
})
app.Get("/site/new", func(c *fiber.Ctx) error {
	// Render index within layouts/main
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("site_new", fiber.Map{
		"title":     "new site",
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	},
		"layouts/main")
})
app.Post("/site/new", func(c *fiber.Ctx) error {
	c.Accepts("application/x-www-form-urlencoded") // "Application/json"

	// todo recevied body is in url format, need to convert to new struct??
	//

	// Render index within layouts/main
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	site := new(models.Site)
	if err := c.BodyParser(site); err != nil {
		log.Warnf("4 %s", err)
		return err
	}

	s, err := site.CreateSite(user.ID, db)
	if err != nil {
		//todo handle error??
		return c.Redirect("/sites")
	}

	_, err = user.AddSite(s, db)
	if err != nil {
		// todo handle error
		return c.Redirect("/sites")
	}

	// todo handle error/success and return to home
	return c.Redirect("/sites")
})
app.Get("/site/:siteid?", func(c *fiber.Ctx) error {
	b, err := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth")
	}
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

	user, err := handler.GetUserFromSession(c, session, db)
	if err != nil {
		return c.Redirect("/auth")
	}

	user.Password = ""

	/*var agentStatList models.AgentStatsList

	stats, err := getAgentStatsForSite(objId, db)
	if err != nil {
		//todo handle error
		//return err
	}
	agentStatList.List = stats

	var hasData = true
	if len(agentStatList.List) == 0 {
		hasData = false
	}

	doc, err := json.Marshal(agentStatList)
	if err != nil {
		log.Errorf("1 %s", err)
	}*/

	// Render index within layouts/main
	// TODO process if they are logged in or not, otherwise send them to registration/login

	return c.Render("site", fiber.Map{
		"title":        "site dashboard",
		"siteSelected": true,
		"siteName":     site.Name,
		"siteId":       site.ID.Hex(),
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email,
		/*"agents":       html.UnescapeString(string(doc)),
		"hasData":      hasData,*/
	},
		"layouts/main")
})
// manage site
