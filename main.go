package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

/*
TODO

api to:
get configuration as agent
post icmp & mtr check as agent


*/

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("NetWatcher Control Server")
	})

	// verify agent / create hash and id
	// GET /api/agent/verify
	app.Get("/api/agent/check/mtr", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("pin"))
		return c.SendString(msg) // => ✋ register
	})

	app.Post("/api/agent/check/mtr", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("pin"))
		return c.SendString(msg) // => ✋ register
	})

	app.Listen(":3000")
}
