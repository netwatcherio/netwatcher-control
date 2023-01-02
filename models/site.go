package models

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Site struct {
	ID              primitive.ObjectID `bson:"_id, omitempty"json:"id"`
	Name            string             `bson:"name"form:"name"json:"name"`
	Members         []SiteMember       `bson:"members"json:"members"`
	CreateTimestamp time.Time          `bson:"create_timestamp"json:"create_timestamp"`
}

type SiteMember struct {
	User primitive.ObjectID `bson:"user"json:"user"`
	Role int                `bson:"role"json:"role"`
	// roles: 0=READ ONLY, 1=READ-WRITE (Create only), 2=ADMIN (Delete Agents), 3=OWNER (Delete Sites)
	// ADMINS can regenerate agent pins
}

type AddSiteMember struct {
	Email string `json:"email"form:"email"`
	Role  int    `json:"role"form:"role"`
}

func (s *Site) CreateSite(owner primitive.ObjectID, db *mongo.Database) (primitive.ObjectID, error) {
	member := SiteMember{
		User: owner,
		Role: 3,
	}

	s.Members = append(s.Members, member)
	//TODO insert into sites list for member
	s.ID = primitive.NewObjectID()
	s.CreateTimestamp = time.Now()

	mar, err := bson.Marshal(s)
	if err != nil {
		log.Errorf("1 %s", err)
		return primitive.ObjectID{}, err
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("2 %s", err)
		return primitive.ObjectID{}, err
	}
	_, err = db.Collection("sites").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return primitive.ObjectID{}, err
	}
	return s.ID, nil
}

// todo when deleting site remove from user document as well

// IsMember check if a user id is a member in the site
func (s *Site) IsMember(id primitive.ObjectID) bool {
	// check if the site contains the member with the provided id
	for _, m := range s.Members {
		if m.User == id {
			return true
		}
	}

	return false
}

// AddMember Add a member to the site then update document
func (s *Site) AddMember(id primitive.ObjectID, role int, db *mongo.Database) (bool, error) {
	// add member with the provided role
	if s.IsMember(id) {
		return false, errors.New("already a member")
	}

	newMember := SiteMember{
		User: id,
		Role: role,
	}

	s.Members = append(s.Members, newMember)
	j, _ := json.Marshal(s.Members)
	log.Warnf("%s", j)

	/*memB, err := bson.Marshal(s.Members)
	if err != nil {
		log.Errorf("69 %s", err)
	}*/

	sites := db.Collection("sites")
	_, err := sites.UpdateOne(
		context.TODO(),
		bson.M{"_id": s.ID},
		bson.D{
			{"$set", bson.D{{"members", s.Members}}},
		},
	)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return true, nil
}

func (s *Site) GetAgents(db *mongo.Database) ([]*Agent, error) {
	// if hash is blank, search for pin matching with blank hash
	// if none exist, return error
	// if match, return new agent, and new hash, then let another function update the hash?
	// if hash is included, search both, and return nil for hash, and false for new if verified
	// if hash is included and none match, return err
	var filter = bson.D{{"site", s.ID}}

	cursor, err := db.Collection("agents").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	//fmt.Println(results)

	if len(results) == 0 {
		return nil, errors.New("no agents match when using id")
	}

	var agent []*Agent
	for i := range results {
		doc, err := bson.Marshal(&results[i])
		if err != nil {
			log.Errorf("1 %s", err)
			return nil, err
		}
		var a *Agent
		err = bson.Unmarshal(doc, &a)
		if err != nil {
			log.Errorf("2 %s", err)
			return nil, err
		}

		agent = append(agent, a)
	}

	return agent, nil
}

func (s *Site) AgentCount(db *mongo.Database) (int, error) {
	var filter = bson.D{{"site", s.ID}}

	count, err := db.Collection("agents").CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

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

	//fmt.Println(results)

	if len(results) > 1 {
		return errors.New("multiple sites match when using id")
	}

	if len(results) == 0 {
		return errors.New("no sites match when using id")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return err
	}

	var site *Site
	err = bson.Unmarshal(doc, &site)
	if err != nil {
		log.Errorf("2 %s", err)
		return err
	}

	s.Name = site.Name
	s.Members = site.Members

	return nil
}

// todo handle if site already exists, or already has the member in it
