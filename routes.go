package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sagostin/netwatcher-agent/agent_models"
	log "github.com/sirupsen/logrus"
)

func LoadApiRoutes(app *fiber.App) {
	app.Post("/v1/agent/update/mtr", func(c *fiber.Ctx) error {
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
			log.Warnf("%s", string(jRespB))
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})

	app.Post("/v1/agent/update/network", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := agent_models.ApiConfigResponse{}
		respB.Response = 200

		log.Warnf("%s", c.Body())

		var data *agent_models.NetworkInfo
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}
		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			log.Warnf("%s", string(jRespB))
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})
	app.Post("/v1/agent/update/speedtest", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := agent_models.ApiConfigResponse{}
		respB.Response = 200

		var data agent_models.SpeedTestInfo

		//log.Warnf("%s", c.Body())

		//fmt.Println(res["json"])

		//log.Infof("%s", string(jMar))
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			log.Warnf("%s", string(jRespB))
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})

	app.Post("/v1/agent/update/icmp", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"
		respB := agent_models.ApiConfigResponse{}
		respB.Response = 200

		var data []*agent_models.IcmpTarget
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			log.Warnf("%s", string(jRespB))
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})
}

func LoadFrontendRoutes(app *fiber.App) {
	app.Get("/404", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("404", fiber.Map{
			"title": "404"})
	})
	app.Get("/", func(c *fiber.Ctx) error {
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
	app.Get("/map", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("map", fiber.Map{
			"title": "map"},
			"layouts/main")
	})
	app.Get("/manage", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("manage", fiber.Map{
			"title": "manage"},
			"layouts/main")
	})
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("dashboard", fiber.Map{
			"title": "dashboard"},
			"layouts/main")
	})
}
