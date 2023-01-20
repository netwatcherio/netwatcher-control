package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/netwatcherio/netwatcher-control/handler"
	"github.com/netwatcherio/netwatcher-control/routes"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	var err error

	runtime.GOMAXPROCS(4)

	log.SetFormatter(&log.TextFormatter{})

	// Load .env
	err = godotenv.Load()
	if err != nil {
		return
	}

	// connect to database
	handler.MongoUri = os.Getenv("MONGO_URI")

	var mongoData *handler.MongoDatastore
	mongoData = handler.NewDatastore(os.Getenv("MAIN_DB"), log.New())

	handleSignals()
	
	router := routes.NewRouter(mongoData.Db)
	router.Init()
	router.Listen(os.Getenv("LISTEN"))
}

func handleSignals() {
	// Signal Termination if using CLI
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGKILL)
	go func() {
		for _ = range signals {
			shutdown()
		}
	}()

}

func shutdown() {
	fmt.Println()
	log.Warnf("%d threads at exit.", runtime.NumGoroutine())
	log.Warn("Shutting down NetWatcher Agent...")
	os.Exit(1)
}
