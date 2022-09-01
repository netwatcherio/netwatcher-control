package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	log "github.com/sirupsen/logrus"
)

func createSession(c *fiber.Ctx, sessions session.Store) error {
	store, err := sessions.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
	}
	defer store.Save()
	return nil
}

func validateSession(c *fiber.Ctx, sessions session.Store) (bool, error) {
	store, err := sessions.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}

	if store.Get("email") == nil {
		return true, errors.New("no email, not logged in")
	}

	defer store.Save()
	return true, nil
}

func loginSession(c *fiber.Ctx, sessions session.Store) (bool, error) {
	store, err := sessions.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}

	if store.Get("email") == nil {
		return true, errors.New("no email, not logged in")
	}

	defer store.Save()
	return true, nil
}
