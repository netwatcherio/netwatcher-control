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

	// Load Auth
	r.auth()
	r.authLogin()
	r.authRegister()

	// Load
}

func LoadFrontendRoutes(app *fiber.App, session *session.Store, db *mongo.Database) {

	// home page

	// dashboard page

	// authentication

	// backend admin TODO
}
