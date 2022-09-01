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

	mainDb, err := MongoConnect(os.Getenv("MAIN_DB"))
	if err != nil {
		log.Fatal("unable to connect to db")
	}

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
		Collection:    "fiber_storage",
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

	LoadApiRoutes(app, store, mainDb)
	LoadFrontendRoutes(app, store, mainDb)

	// Listen website
	app.Listen(os.Getenv("LISTEN"))
}

func shutdown() {
	log.Fatal("Shutting down NetWatcher Agent...")
}
