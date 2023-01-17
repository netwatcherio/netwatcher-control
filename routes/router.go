package routes

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/netwatcherio/netwatcher-control/handler/agent"
	"github.com/netwatcherio/netwatcher-control/workers"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type Router struct {
	App     *fiber.App
	Session *session.Store
	DB      *mongo.Database
}

func NewRouter(store *mongo.Database) *Router {
	router := Router{
		App: fiber.New(),
		DB:  store,
	}
	return &router
}

func (r *Router) Init() {
	checkCreateWorker := make(chan agent.Data)

	if os.Getenv("DEBUG") != "" {
		log.Warning("Cross Origin requests allowed (ENV::DEBUG)")
		r.App.Use(cors.New())
	}

	log.Info("AUTH")
	// Initialize the auth routes before using JWT
	r.login()
	r.register()

	rng := rand.Reader
	privateKey, err := rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}

	// JWT Middleware
	r.App.Use(jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    privateKey.Public(),
	}))

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

	workers.CreateCheckWorker(checkCreateWorker, r.DB)
}

func (r *Router) Listen(host string) {
	err := r.App.Listen(host)
	if err != nil {
		return
	}
}
