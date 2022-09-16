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
		return c.Render("home", fiber.Map{
			"title":     "home",
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		},
			"layouts/main")
	})
	// dashboard page

	app.Get("/agent/:agent?", func(c *fiber.Ctx) error {
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

		icmpD, err := getIcmpData(objId, time.Minute*30, db)
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
		return c.Render("agent", fiber.Map{
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
		if c.Params("siteid") == "" {
			return c.Redirect("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("siteid"))
		if err != nil {
			return c.Redirect("/home")
		}

		site, err := getSite(objId, db)
		if err != nil {
			return nil
		}

		var agentStatList control_models.AgentStatsList

		stats, err := getAgentStats(objId, db)
		if err != nil {
			return err
		}
		agentStatList.List = stats

		doc, err := json.Marshal(agentStatList)
		if err != nil {
			log.Errorf("1 %s", err)
		}

		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		log.Errorf("%s", string(doc))
		return c.Render("agents", fiber.Map{
			"title": "agents", "siteSelected": true, "siteName": site.Name, "siteId": site.ID.Hex(), "agents": html.UnescapeString(string(doc))},
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

		var sitesList struct {
			Sites []*control_models.Site `json:"sites"`
		}

		for _, sid := range user.Sites {
			site, err := getSite(sid, db)
			if err != nil {
				// todo display error instead of redirecting
				log.Errorf("%s", err)
			}

			sitesList.Sites = append(sitesList.Sites, site)
		}

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
			"sites":     siteJs,
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

		var sitesList struct {
			Sites []*control_models.Site `json:"sites"`
		}

		for _, sid := range user.Sites {
			site, err := getSite(sid, db)
			if err != nil {
				// todo display error instead of redirecting
				log.Errorf("%s", err)
			}

			sitesList.Sites = append(sitesList.Sites, site)
		}

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
			"sites":     siteJs,
		},
			"layouts/main")
	})
	app.Post("/site/new", func(c *fiber.Ctx) error {
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
			return nil
		}

		var agentStatList control_models.AgentStatsList

		stats, err := getAgentStats(objId, db)
		if err != nil {
			return err
		}
		agentStatList.List = stats

		doc, err := json.Marshal(agentStatList)
		if err != nil {
			log.Errorf("1 %s", err)
		}

		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		log.Errorf("%s", string(doc))
		return c.Render("site", fiber.Map{
			"title": "site dashboard", "siteSelected": true, "siteName": site.Name, "siteId": site.ID.Hex(), "agents": html.UnescapeString(string(doc))},
			"layouts/main")
	})
	// manage site
	app.Get("/site/manage", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("manage", fiber.Map{
			"title": "manage"},
			"layouts/main")
	})

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
