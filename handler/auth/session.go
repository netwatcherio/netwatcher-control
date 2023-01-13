package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	UserID    primitive.ObjectID `json:"user_id"bson:"user_id"`
	SessionID primitive.ObjectID `json:"session_id"bson:"_id"`
	Expiry    time.Time          `json:"expiry"bson:"expiry"`
}

func (s *Session) Check() error {
	return nil
}
