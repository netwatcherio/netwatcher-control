package routes

import "github.com/gofiber/fiber/v2"

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) alerts() {
	// alerts
	r.App.Get("/alerts", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("alerts", fiber.Map{
			"title": "alerts"},
			"layouts/main")
	})
}
