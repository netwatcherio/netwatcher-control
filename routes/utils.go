package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) utils() {
	r.App.Get("/404", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("404", fiber.Map{
			"title": "404"})
	})

	r.App.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Redirect("/sites")
	})
	r.App.Get("/home", func(c *fiber.Ctx) error {
		_, err := validateUser(r, c)
		if err != nil {
			return err
		}
		return c.Redirect("/sites")
	})
}
