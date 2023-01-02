package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	App     *fiber.App
	Session *session.Store
	DB      *mongo.Database
}

func (r *Router) Init() {
	log.Info("Loading routes for:")

	r.apiGetConfig()
	log.Info("API")

	r.utils()
	log.Info("UTILS")

	r.agent()
	r.agents()
	r.agentNew()
	r.agentInstall()
	log.Info("AGENTS")

	r.alerts()
	log.Info("ALERTS")

	r.site()
	r.sites()
	r.siteNew()
	r.siteMembers()
	r.siteAddMember()
	log.Info("SITES")

	r.auth()
	r.authLogin()
	r.authRegister()
	r.authLogout()
	log.Info("AUTH")
}
