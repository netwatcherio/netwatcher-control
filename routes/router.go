package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/netwatcherio/netwatcher-control/handler/agent"
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
	checkCreateWorker := make(chan agent.Data)

	log.Info("Loading routes for:")

	r.apiGetConfig()
	r.apiDataPush(checkCreateWorker)
	log.Info("API")

	r.getCheck()
	r.getCheckData()
	r.checkNew()
	log.Info("CHECKS")

	r.agentNew()
	r.getAgent()
	r.getGeneralAgentStats()
	log.Info("AGENTS")

	r.addSiteMember()
	r.getSite()
	r.getSites()
	r.newSite()
	log.Info("SITES")

	r.login()
	r.register()
	log.Info("AUTH")

	workers.CreateCheckWorker(checkCreateWorker, r.DB)
}
