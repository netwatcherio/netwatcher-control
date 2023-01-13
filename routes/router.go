package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/netwatcherio/netwatcher-control/handler"
	"github.com/netwatcherio/netwatcher-control/workers"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	App     *fiber.App
	Session *session.Store
	DB      *mongo.Database
}

func (r *Router) Init() {
	checkCreateWorker := make(chan handler.CheckData)

	log.Info("Loading routes for:")

	r.apiGetConfig()
	r.apiDataPush(checkCreateWorker)
	log.Info("API")

	r.check()
	r.checkNew()
	log.Info("CHECKS")

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

	r.login()
	/*r.authLogin()
	r.authRegister()
	r.authLogout()*/
	log.Info("AUTH")

	workers.CreateCheckWorker(checkCreateWorker, r.DB)
}
