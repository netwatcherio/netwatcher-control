package workers

import (
	"github.com/netwatcherio/netwatcher-control/handler/agent"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCheckWorker(c chan agent.Data, db *mongo.Database) {
	go func(cl chan agent.Data) {
		log.Info("Starting check data creation worker...")
		for {
			data := <-cl

			err := data.Create(db)
			if err != nil {
				log.Error(err)
			}

			if data.Type == agent.CtSpeedTest {
				agentC := agent.Check{ID: data.CheckID}
				_, err := agentC.Get(db)
				if err != nil {
					log.Error(err)
					return
				}
				agentC.Pending = false
				err = agentC.Update(db)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}(c)
}
