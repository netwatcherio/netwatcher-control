package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginSession(c *fiber.Ctx, session *session.Store, db *mongo.Database, id primitive.ObjectID) (bool, error) {
	store, err := session.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}

	store.Set("id", id)
	defer store.Save()
	return true, nil
}

func ValidateSession(c *fiber.Ctx, sessions *session.Store, db *mongo.Database) (bool, error) {
	store, err := sessions.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}

	if store.Get("id") == nil {
		return true, errors.New("no email, not logged in")
	}
	// todo does this work?? do i need to compare against db??
	return true, nil

	/*userId, err := primitive.ObjectIDFromHex(store.Get("id").(string))
	if err != nil {
		return false, err
	}

	user := control_models.User{ID: userId}
	b, err := user.UserExistsID(db)
	if err != nil {
		return false, err
	}

	defer store.Save()

	if b {
		return true, nil
	}*/
	return false, nil
}

func LogoutSession(c *fiber.Ctx, sessions *session.Store) (bool, error) {
	store, err := sessions.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}

	store.Set("id", nil)

	defer store.Save()
	return true, nil
}
