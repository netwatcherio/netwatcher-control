package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/netwatcherio/netwatcher-control/control_models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"html"
	"math"
	"strings"
	"time"
)

func LoadApiRoutes(app *fiber.App, session *session.Store, db *mongo.Database) {
	app.Post("/v1/agent/update/mtr", func(c *fiber.Ctx) error {
		var str = apiUpdateMtr(c, db)
		if str != "" {
			return c.SendString(str)
		}
		return c.SendString("Something went wrong...") // => ✋
	})

	app.Post("/v1/agent/update/network", func(c *fiber.Ctx) error {
		var str = apiUpdateNetwork(c, db)
		if str != "" {
			return c.SendString(str)
		}

		return c.SendString("Something went wrong...") // => ✋
	})
	app.Post("/v1/agent/update/speedtest", func(c *fiber.Ctx) error {
		var str = apiUpdateSpeedTest(c, db)
		if str != "" {
			return c.SendString(str)
		}

		return c.SendString("Something went wrong...") // => ✋
	})

	app.Post("/v1/agent/update/icmp", func(c *fiber.Ctx) error {
		var str = apiUpdateIcmp(c, db)
		if str != "" {
			return c.SendString(str)
		}

		return c.SendString("Something went wrong...") // => ✋
	})

	app.Get("/v1/agent/config/:pin/:hash?", func(c *fiber.Ctx) error {
		var str = apiGetConfig(c, db)
		if str != "" {
			return c.SendString(str)
		}

		return c.SendString("Something went wrong...") // => ✋
	})
}

func LoadFrontendRoutes(app *fiber.App, session *session.Store, db *mongo.Database) {
	app.Get("/404", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("404", fiber.Map{
			"title": "404"})
	})
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Redirect("/home")
	})

	// home page
	app.Get("/home", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		// convert to json for testing
		_, err = json.Marshal(user.Sites)
		if err != nil {
			// todo handle properly
			return c.Redirect("/auth")
		}

		//TODO get list of sites based on sites on user

		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		/*return c.Render("home", fiber.Map{
			"title":     "home",
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		},
			"layouts/main")*/

		return c.Redirect("/sites")
	})
	// dashboard page

	app.Get("/agent/new/:siteid?", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			// todo handle error
			//return nil
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
	app.Post("/agent/new/:siteid?", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		// Render index within layouts/main
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		if c.Params("siteid") == "" {
			return c.Redirect("/agents")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.Redirect("/agents")
		}

		site, err := getSite(objId, db)
		if err != nil {
			// todo handle error
			//return nil
		}

		cAgent := new(control_models.CreateAgent)
		if err := c.BodyParser(cAgent); err != nil {
			log.Warnf("4 %s", err)
			return err
		}

		icmpTargets := strings.Split(cAgent.IcmpTargets, ",")
		mtrTargets := strings.Split(cAgent.MtrTargets, ",")

		agentId, err := CreateAgent(cAgent.Name, icmpTargets, mtrTargets, site.ID, db)
		if err != nil {
			//todo handle error??
			return c.Redirect("/agents")
		}

		// todo handle error/success and return to home
		return c.Redirect("/agent/install/" + agentId.String())
	})
	app.Get("/agent/install/:agentid", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		if c.Params("agentid") == "" {
			return c.Redirect("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("agentid"))
		if err != nil {
			return c.Redirect("/home")
		}

		agent, err := getAgent(objId, db)
		if err != nil {
			// todo handle error
			//return nil
			return c.Redirect("/agents")
		}

		site, err := getSite(agent.Site, db)
		if err != nil {
			// todo handle error
			return c.Redirect("/agents")
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

	app.Get("/agent/:agent?", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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

		agent, err := getAgent(objId, db)
		if err != nil {
			log.Errorf("1 %s", err)
		}

		site, err := getSite(agent.Site, db)
		if err != nil {
			log.Errorf("12 %s", err)
		}

		marshal, err := json.Marshal(agent)
		if err != nil {
			log.Errorf("13 %s", err)
		}

		networkData, err := getLatestNetworkData(agent.ID, db)
		if err != nil {
			log.Errorf("14 %s", err)
		}

		mtrData, err := getLatestMtrData(agent.ID, 5, db)
		if err != nil {
			log.Errorf("15 %s", err)
		}

		marshalMtr, err := json.Marshal(mtrData)
		if err != nil {
			log.Errorf("16 %s", err)
		}

		speedData, err := getLatestSpeedtests(agent.ID, 5, db)
		if err != nil {
			log.Errorf("15 %s", err)
		}

		marshalSpeed, err := json.Marshal(speedData)
		if err != nil {
			log.Errorf("16 %s", err)
		}

		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("agent", fiber.Map{
			"title":            agent.Name,
			"siteSelected":     true,
			"siteName":         site.Name,
			"siteId":           site.ID.Hex(),
			"agents":           html.UnescapeString(string(marshal)),
			"mtr":              html.UnescapeString(string(marshalMtr)),
			"speed":            html.UnescapeString(string(marshalSpeed)),
			"publicAddress":    networkData.Data.PublicAddress,
			"localAddress":     networkData.Data.LocalAddress,
			"defaultGateway":   networkData.Data.DefaultGateway,
			"internetProvider": networkData.Data.InternetProvider,
			"uploadSpeed":      math.Round(speedData[0].Data.ULSpeed),
			"downloadSpeed":    math.Round(speedData[0].Data.DLSpeed),
			"firstName":        user.FirstName,
			"lastName":         user.LastName,
			"email":            user.Email,
			//"icmpMetrics":  html.UnescapeString(string(j)),
		},
			"layouts/main")
	})
	app.Get("/icmp/:agent?", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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

		agent, err := getAgent(objId, db)
		if err != nil {
			log.Errorf("1 %s", err)
			return err
		}

		site, err := getSite(agent.Site, db)
		if err != nil {
			log.Errorf("12 %s", err)
			return err
		}

		doc, err := json.Marshal(agent)
		if err != nil {
			log.Errorf("13 %s", err)
			return err
		}

		icmpD, err := getIcmpData(objId, time.Hour*24, db)
		if err != nil {
			return err
		}

		for n := range icmpD {
			for n2 := range icmpD[n].Data {
				icmpD[n].Data[n2].Result.Data = nil
			}
		}

		j, err := json.Marshal(icmpD)
		if err != nil {
			log.Errorf("%s", err)
			return err
		}

		log.Errorf("%s", j)

		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("icmp", fiber.Map{
			"title":        agent.Name,
			"siteSelected": true,
			"siteName":     site.Name,
			"siteId":       site.ID.Hex(),
			"agents":       html.UnescapeString(string(doc)),
			"icmpMetrics":  html.UnescapeString(string(j)),
		},
			"layouts/main")
	})
	app.Get("/agents/:siteid?", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			// todo handle error
			//return nil
		}

		var agentStatList control_models.AgentStatsList

		stats, err := getAgentStats(objId, db)
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
		}

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
			"email":        user.Email,
			"agents":       html.UnescapeString(string(doc)),
			"hasAgents":    hasAgents},
			"layouts/main")
	})

	app.Get("/traceroutes/:agent?", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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

		agent, err := getAgent(objId, db)
		if err != nil {
			log.Errorf("1 %s", err)
			return err
		}

		site, err := getSite(agent.Site, db)
		if err != nil {
			log.Errorf("12 %s", err)
			return err
		}
		var agentStatList control_models.AgentStatsList

		stats, err := getAgentStats(objId, db)
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
		}

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
			"email":        user.Email,
			"agents":       html.UnescapeString(string(doc)),
			"hasAgents":    hasAgents},
			"layouts/main")
	})

	app.Get("/alerts/:siteid?", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			// todo handle error
			//return nil
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

	// authentication
	app.Get("/auth/register", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if b {
			return c.Redirect("/home")
		}
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("auth", fiber.Map{
			"title": "auth", "login": false})
	})
	app.Get("/auth/login", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if b {
			return c.Redirect("/home")
		}
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("auth", fiber.Map{
			"title": "auth", "login": true})
	})
	app.Get("/auth", func(c *fiber.Ctx) error {
		b, _ := ValidateSession(c, session, db)
		if b {
			return c.Redirect("/home")
		}
		return c.Redirect("/auth/login")
	})
	app.Post("/auth/register", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		registerUser := new(control_models.RegisterUser)
		if err := c.BodyParser(registerUser); err != nil {
			log.Warnf("%s", err)
			return err
		}

		if registerUser.Password != registerUser.PasswordConfirm {
			//todo handle error and show on auth page using sessions??
			return c.Redirect("/auth/register")
		}

		pwd, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 15)
		if err != nil {
			log.Errorf("%s", err)
			return c.Redirect("/auth/login")
		}

		user := control_models.User{
			ID:        primitive.NewObjectID(),
			Email:     registerUser.Email,
			FirstName: registerUser.FirstName,
			LastName:  registerUser.LastName,
			Admin:     false,
			Password:  string(pwd),
			Sites:     nil,
			Verified:  false,
		}

		ucb, err2 := user.Create(db)
		if err2 != nil || !ucb {
			log.Infof("%s", "error creating user")
			return c.Redirect("/auth/register")
		}

		//todo handle success and send to login page
		return c.Redirect("/auth/login")
	})
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		loginUser := new(control_models.LoginUser)
		if err := c.BodyParser(loginUser); err != nil {
			log.Warnf("4 %s", err)
			return err
		}

		user := control_models.User{Email: loginUser.Email}

		// get user from email
		usr, err2 := user.GetUserFromEmail(db)
		if err2 != nil {
			log.Warnf("3 %s", err2)
			return c.Redirect("/auth/login")
		}

		err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(loginUser.Password))
		if err != nil {
			log.Errorf("%s", err)
			return c.Redirect("/auth/login")
		}

		// create token session
		b, err := LoginSession(c, session, db, usr.ID)
		if err != nil || !b {
			log.Warnf("5 %s, 2 %b", err, b)
			return c.Redirect("/auth/login")
		}
		// todo handle success and return to home
		return c.Redirect("/home")
	})
	app.Get("/logout", func(c *fiber.Ctx) error {
		LogoutSession(c, session)

		return c.Redirect("/auth")
	})

	app.Get("/site/:siteid?/members/add", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		b, _ := ValidateSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			//todo handle error
			//return nil
		}

		user, err := GetUserFromSession(c, session, db)
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
		b, _ := ValidateSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			//todo handle error
			//return nil
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		newMember := new(control_models.AddSiteMember)
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

		var usrTmp = control_models.User{Email: newMember.Email}

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
		b, _ := ValidateSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			//todo handle error
			//return nil
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		var siteMembers []control_models.SiteMember
		for _, mem := range site.Members {
			siteMembers = append(siteMembers, mem)
		}

		var siteUsers []*control_models.User
		for _, usr := range siteMembers {
			c2 := &control_models.User{ID: usr.User}
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
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		type AgentCountInfo struct {
			SiteID primitive.ObjectID `json:"site_id"`
			Count  int                `json:"count"`
		}

		var sitesList struct {
			Sites          []*control_models.Site `json:"sites"`
			AgentCountInfo []AgentCountInfo       `json:"agentCountInfo"`
		}

		for _, sid := range user.Sites {
			site, err := getSite(sid, db)
			if err != nil {
				// todo display error instead of redirecting
				log.Errorf("%s", err)
			}

			count, err := getAgentCount(site.ID, db)
			if err != nil {
				//todo handle error
			}

			tempCount := AgentCountInfo{
				SiteID: site.ID,
				Count:  count,
			}

			sitesList.Sites = append(sitesList.Sites, site)
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
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
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
		b, _ := ValidateSession(c, session, db)
		if !b {
			return c.Redirect("/auth/login")
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		site := new(control_models.Site)
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
		b, err := ValidateSession(c, session, db)
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

		site, err := getSite(objId, db)
		if err != nil {
			//todo handle error
			//return nil
		}

		user, err := GetUserFromSession(c, session, db)
		if err != nil {
			return c.Redirect("/auth")
		}

		user.Password = ""

		var agentStatList control_models.AgentStatsList

		stats, err := getAgentStats(objId, db)
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
		}

		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		log.Errorf("%s", string(doc))
		return c.Render("site", fiber.Map{
			"title":        "site dashboard",
			"siteSelected": true,
			"siteName":     site.Name,
			"siteId":       site.ID.Hex(),
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"agents":       html.UnescapeString(string(doc)),
			"hasData":      hasData,
		},
			"layouts/main")
	})
	// manage site

	// alerts
	app.Get("/alerts", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("alerts", fiber.Map{
			"title": "alerts"},
			"layouts/main")
	})
	// profile
	app.Get("/profile", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("profile", fiber.Map{
			"title": "profile"},
			"layouts/main")
	})

	// backend admin TODO
}
