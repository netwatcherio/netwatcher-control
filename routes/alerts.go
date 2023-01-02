package routes

import "github.com/gofiber/fiber/v2"

// alerts
app.Get("/alerts", func(c *fiber.Ctx) error {
	// Render index within layouts/main
	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("alerts", fiber.Map{
		"title": "alerts"},
		"layouts/main")
})
