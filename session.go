package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	log "github.com/sirupsen/logrus"
)

func CreateSession(c *fiber.Ctx, sessions session.Store) (bool, error) {
	store, err := sessions.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
	}
	defer store.Save()
	return true, nil
}

func ValidateSession(c *fiber.Ctx) (bool, error) {
	s := session.New()
	store, err := s.Get(c) // get/create new session
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}
	defer store.Save()
	return true, nil
}
