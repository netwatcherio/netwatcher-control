package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	Debug = false
)

func main() {
	var err error
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.TextFormatter{})

	// Load .env
	godotenv.Load()
	if os.Getenv("DEBUG") == "true" {
		Debug = true
	}

	// connect to database
	mongoUri = os.Getenv("MONGO_URI")

	var mongoData *MongoDatastore
	mongoData = NewDatastore(os.Getenv("MAIN_DB"), log.New())

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

	// Initialize custom config
	sessionStore := mongodb.New(mongodb.Config{
		ConnectionURI: mongoUri,
		Collection:    os.Getenv("SESSIONS_COLLECTION"),
		Database:      os.Getenv("SESSIONS_DB"),
		Reset:         false,
	})

	store := session.New(session.Config{
		Storage: sessionStore,
	})

	// Create a new engine
	engine := html.New("./views", ".html")
	// Load Fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	// Layout defines the variable name that is used to yield templates within layouts
	engine.Layout("embed") // Optional. Default: "embed"

	// Delims sets the action delimiters to the specified strings
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters

	// Public Files
	app.Static("/", "./public")

	//createAgent(mongoData.db)
	//createSite(mongoData.db)

	LoadApiRoutes(app, store, mongoData.db)
	LoadFrontendRoutes(app, store, mongoData.db)

	// Listen website
	app.Listen(os.Getenv("LISTEN"))
}

func shutdown() {
	log.Fatal("Shutting down NetWatcher Agent...")
}
