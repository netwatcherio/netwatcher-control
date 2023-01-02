package routes

import "github.com/gofiber/fiber/v2"

// profile
app.Get("/profile", func(c *fiber.Ctx) error {
	// Render index within layouts/main
	// TODO process if they are logged in or not, otherwise send them to registration/login
	return c.Render("profile", fiber.Map{
		"title": "profile"},
		"layouts/main")
})