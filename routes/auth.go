package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
	log "github.com/sirupsen/logrus"
)

func (r *Router) login() {
	r.App.Post("/auth/login", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"

		var l auth.Login
		err := json.Unmarshal(c.Body(), &l)
		if err != nil {
			log.Error(err)
		}
		t, err := l.Login(r.DB)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"token": t})
	})
}

func (r *Router) register() {
	r.App.Post("/auth/register", func(c *fiber.Ctx) error {
		c.Accepts("Application/json") // "Application/json"

		var reg auth.Register
		err := json.Unmarshal(c.Body(), &reg)
		if err != nil {
			log.Error(err)
		}
		t, err := reg.Register(r.DB)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"token": t})
	})
}
