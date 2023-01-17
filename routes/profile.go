package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
)

// getProfile
func (r *Router) getProfile() {
	r.App.Get("/profile", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"

		t := c.Locals("user").(*jwt.Token)

		user, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		user.Password = ""

		marshal, err := json.Marshal(*user)
		if err != nil {
			return c.JSON(err)
		}

		return c.Send(marshal)
	})
}

//
//import (
//	"github.com/gofiber/fiber/v2"
//)
//
//// TODO authenticate & verify that the user is infact apart of the site etc.
//
//// profile
//app.Get("/profile", func(c *fiber.Ctx) error {
//	// Render index within layouts/main
//	// TODO process if they are logged in or not, otherwise send them to registration/login
//	return c.Render("profile", fiber.Map{
//		"title": "profile"},
//		"layouts/main")
//})
