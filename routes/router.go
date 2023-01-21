package routes

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
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
func secretKey() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		// Always check the signing method
		if t.Method.Alg() != jwtware.HS256 {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		signingKey := os.Getenv("KEY")

		return []byte(signingKey), nil
	}
}

var privateKey *rsa.PrivateKey

func (r *Router) Init() {
	checkCreateWorker := make(chan agent.Data)

	if os.Getenv("DEBUG") != "" {
		log.Warning("Cross Origin requests allowed (ENV::DEBUG)")
		r.App.Use(cors.New())
	}

	var err error
	rng := rand.Reader
	privateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}

	log.Info("AUTH")
	// Initialize the auth routes before using JWT
	r.login(privateKey)
	r.register(privateKey)
	r.apiGetConfig()
	r.apiDataPush(checkCreateWorker)
	
	workers.CreateCheckWorker(checkCreateWorker, r.DB)

	log.Info("API")
	// JWT Middleware
	r.App.Use(jwtware.New(jwtware.Config{
		KeyFunc: secretKey(),
	}))

	log.Info("Loading routes for:")

	r.getCheck()
	r.getChecks()
	r.getCheckData()
	r.checkNew()
	r.deleteCheck()
	log.Info("CHECKS")

	r.agentNew()
	r.getAgents()
	r.getAgent()
	r.getGeneralAgentStats()
	r.deleteAgent()
	log.Info("AGENTS")

	r.addSiteMember()
	r.getSite()
	r.getSites()
	r.newSite()
	r.deleteSite()
	log.Info("SITES")

	r.getProfile()

}

func (r *Router) Listen(host string) {
	err := r.App.Listen(host)
	if err != nil {
		return
	}
}
