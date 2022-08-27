package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

/*
TODO

api to:
get configuration as agent
post icmp & mtr check as agent


*/

func main() {
	var err error
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.TextFormatter{})

	godotenv.Load()

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("NetWatcher Control Server")
	})

	app.Post("/api/agent/update/icmp", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("data"))

		err := json.Unmarshal(msg)

		log.Infof(msg)
		return c.SendString(msg) // => ✋ register
	})

	app.Listen(":3000")
}

func shutdown() {
	log.Fatal("Shutting down NetWatcher Agent...")
}
