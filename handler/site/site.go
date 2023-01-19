package site

import (
	"context"
	"errors"
	"github.com/netwatcherio/netwatcher-control/handler/agent"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Site struct {
	ID              primitive.ObjectID `bson:"_id, omitempty"json:"id"`
	Name            string             `bson:"name"form:"name"json:"name"`
	Members         []Member           `bson:"members"json:"members"`
	CreateTimestamp time.Time          `bson:"create_timestamp"json:"create_timestamp"`
}

func (s *Site) Create(owner primitive.ObjectID, db *mongo.Database) error {
	member := Member{
		User: owner,
		Role: SrADMIN,
	}

	s.Members = append(s.Members, member)
	s.ID = primitive.NewObjectID()
	s.CreateTimestamp = time.Now()

	mar, err := bson.Marshal(s)
	if err != nil {
		return errors.New("unable to marshal bson data")
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		return errors.New("unable to marshal site data")
	}
	_, err = db.Collection("sites").InsertOne(context.TODO(), b)
	if err != nil {
		return errors.New("unable to create site")
	}

	u := auth.User{ID: member.User}
	usr, err := u.FromID(db)
	if err != nil {
		return errors.New("unable to get user from id")
	}
	u = *usr
	err = u.AddSite(s.ID, db)
	if err != nil {
		return errors.New("unable to add site to user")
	}

	return nil
}

// todo when deleting site remove from user document as well

func (s *Site) GetAgents(db *mongo.Database) ([]*agent.Agent, error) {
	var filter = bson.D{{"site", s.ID}}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, errors.New("unable to search database for agents")
	}

	if len(results) == 0 {
		return nil, errors.New("no agents match when using id")
	}

	var agents []*agent.Agent
	for i := range results {
		doc, err := bson.Marshal(&results[i])
		if err != nil {
			return nil, err
		}
		var a *agent.Agent
		err = bson.Unmarshal(doc, &a)
		if err != nil {
			return nil, err
		}

		agents = append(agents, a)
	}

	return agents, nil
}

// AgentCount returns count of agents for a site, or an error if its not successful
func (s *Site) AgentCount(db *mongo.Database) (int, error) {
	var filter = bson.D{{"site", s.ID}}

	count, err := db.Collection("agents").CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Get a site from the provided ID
func (s *Site) Get(db *mongo.Database) error {
	var filter = bson.D{{"_id", s.ID}}

	cursor, err := db.Collection("sites").Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return err
	}

	if len(results) > 1 {
		return errors.New("multiple sites match when using id")
	}

	if len(results) == 0 {
		return errors.New("no sites match when using id")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		return err
	}

	var site Site
	err = bson.Unmarshal(doc, &site)
	if err != nil {
		return err
	}

	s.Name = site.Name
	s.Members = site.Members
	s.CreateTimestamp = site.CreateTimestamp

	return nil
}

// Delete data based on provided agent ID in checkData struct
func (s *Site) Delete(db *mongo.Database) error {
	// filter based on check ID
	var filter = bson.D{{"_id", s.ID}}

	_, err := db.Collection("sites").DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) SiteStats(db *mongo.Database) ([]*agent.Stats, error) {
	var agentStats []*agent.Stats

	agents, err := s.GetAgents(db)
	if err != nil {
		return nil, err
	}
	for _, a := range agents {
		stats, err := a.GetLatestStats(db)
		if err != nil {
			log.Error(err)
		}
		agentStats = append(agentStats, stats)
	}

	return agentStats, nil
}
