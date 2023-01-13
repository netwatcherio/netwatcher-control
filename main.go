package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
	"github.com/netwatcherio/netwatcher-control/handler"
	"github.com/netwatcherio/netwatcher-control/routes"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	Debug = false
)

var (
	// Obviously, this is just a test example. Do not do this in production.
	// In production, you would have the private key and public key pair generated
	// in advance. NEVER add a private key to any GitHub repo.
	privateKey *rsa.PrivateKey
)

func main() {
	var err error
	if err != nil {
		log.Fatal(err)
	}

	runtime.GOMAXPROCS(4)

	log.SetFormatter(&log.TextFormatter{})

	// Load .env
	godotenv.Load()
	if os.Getenv("DEBUG") == "true" {
		Debug = true
	}

	// connect to database
	handler.MongoUri = os.Getenv("MONGO_URI")

	var mongoData *handler.MongoDatastore
	mongoData = handler.NewDatastore(os.Getenv("MAIN_DB"), log.New())

	// Signal Termination if using CLI
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGKILL)
	go func() {
		s := <-signals
		log.Fatal("Received Signal: %s", s)
		shutdown()
		os.Exit(1)
	}()

	app := fiber.New()

	// Just as a demo, generate a new private/public key pair on each run. See note above.
	rng := rand.Reader
	privateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    privateKey.Public(),
	}))

	//createAgent(mongoData.db)
	//createSite(mongoData.db)

	router := routes.Router{
		App: app,
		DB:  mongoData.Db,
	}

	router.Init()

	// Listen website
	app.Listen(os.Getenv("LISTEN"))
}

func shutdown() {
	log.Fatal("Shutting down NetWatcher Agent...")
}
