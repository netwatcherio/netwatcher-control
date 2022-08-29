package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sagostin/netwatcher-agent/agent_models"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	Debug = false
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
	if os.Getenv("DEBUG") == "true" {
		Debug = true
	}

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

	app.Post("/v1/agent/update/mtr", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"
		respB := agent_models.ApiConfigResponse{}
		respB.Response = 200

		log.Warnf("%s", c.Body())

		var icmpData []*agent_models.MtrTarget

		//fmt.Println(res["json"])

		//log.Infof("%s", string(jMar))
		err = json.Unmarshal(c.Body(), &icmpData)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			log.Warnf("%s", string(jRespB))
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})

	app.Post("/v1/agent/update/icmp", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"
		respB := agent_models.ApiConfigResponse{}
		respB.Response = 200

		var icmpData []*agent_models.IcmpTarget

		//fmt.Println(res["json"])

		//log.Infof("%s", string(jMar))
		err = json.Unmarshal(c.Body(), &icmpData)
		if err != nil {
			log.Errorf("2 %s", err)
			respB.Response = 500
		}

		jRespB, err := json.Marshal(respB)
		if err != nil {
			log.Errorf("3 Unable to marshal API response.")
		} else {
			log.Warnf("%s", string(jRespB))
			return c.SendString(string(jRespB)) // => ✋ good
		}

		return c.SendString("Something went wrong...") // => ✋
	})

	app.Listen(os.Getenv("LISTEN"))
}

func shutdown() {
	log.Fatal("Shutting down NetWatcher Agent...")
}
