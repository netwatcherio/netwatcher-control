package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	App     *fiber.App
	Session *session.Store
	DB      *mongo.Database
}

func (r *Router) Init() {
	r.utils()

	r.agent()
	r.agents()
	r.agentNew()
	r.agentInstall()

	r.alerts()

	// Load Auth
	r.auth()
	r.authLogin()
	r.authRegister()

	// Load
}
