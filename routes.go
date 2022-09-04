package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sagostin/netwatcher-agent/agent_models"
	"github.com/sagostin/netwatcher-control/control_models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html"
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
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("home", fiber.Map{
			"title": "home"},
			"layouts/main")
	})
	// dashboard page
	app.Get("/dashboard/:siteid?", func(c *fiber.Ctx) error {
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
		return c.Render("dashboard", fiber.Map{
			"title": "dashboard", "siteSelected": true, "siteName": site.Name, "agents": html.UnescapeString(string(doc))},
			"layouts/main")
	})

	// authentication
	app.Get("/auth", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("auth", fiber.Map{
			"title": "auth"})
	})
	app.Post("/auth", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := agent_models.ApiConfigResponse{}
		respB.Response = 200

		var data []*agent_models.MtrTarget
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}
		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})
	app.Get("/logout", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("index", fiber.Map{
			"title": "home"},
			"layouts/main")
	})

	app.Get("/agents", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("agents", fiber.Map{
			"title": "agents"},
			"layouts/main")
	})
	app.Get("/sites", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("sites", fiber.Map{
			"title": "sites"},
			"layouts/main")
	})
	app.Get("/alerts", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("alerts", fiber.Map{
			"title": "alerts"},
			"layouts/main")
	})
	app.Get("/profile", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("profile", fiber.Map{
			"title": "profile"},
			"layouts/main")
	})
	app.Get("/manage", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("manage", fiber.Map{
			"title": "manage"},
			"layouts/main")
	})
}
