package main

import (
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

	router := routes.NewRouter(mongoData.Db)
	router.Init()
	router.Listen(os.Getenv("LISTEN"))
}

func shutdown() {
	log.Fatal("Shutting down NetWatcher Agent...")
}
