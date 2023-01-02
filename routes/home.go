package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"netwatcher-control/handler"
)

app.Get("/home", func(c *fiber.Ctx) error {
	b, _ := handler.ValidateSession(c, session, db)
	if !b {
		return c.Redirect("/auth/login")
	}

	user, err := handler.GetUserFromSession(c, session, db)
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