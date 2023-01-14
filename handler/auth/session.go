package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Session struct {
	UserID primitive.ObjectID `json:"user_id"bson:"user_id"`
	ID     primitive.ObjectID `json:"session_id"bson:"_id"`
	Expiry time.Time          `json:"expiry"bson:"expiry"`
}

// Create a session from user id, and include expiry, return error if fails
func (s *Session) Create(db *mongo.Database) error {
	s.ID = primitive.NewObjectID()
	s.Expiry = time.Now().Add(time.Hour * 24)

	if (s.UserID == primitive.ObjectID{}) {
		return errors.New("invalid userid used to create session")
	}

	mar, err := bson.Marshal(s)
	if err != nil {
		return errors.New("something went wrong marshalling session struct")
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		return errors.New("something went wrong marshalling session struct")
	}

	_, err = db.Collection("sessions").InsertOne(context.TODO(), b)
	if err != nil {
		return errors.New("something went wrong marshalling session struct")
	}

	return nil
}

// FromID returns a user if it finds a matching user with the provided ID
func (s *Session) FromID(db *mongo.Database) (*Session, error) {
	var filter = bson.D{{"_id", s.ID}}
	cursor, err := db.Collection("sessions").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) < 1 {
		return nil, errors.New("no session found")
	}

	if len(results) > 1 {
		return nil, errors.New("multiple sessions found")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	var session *Session
	err = bson.Unmarshal(doc, &session)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, errors.New("something went wrong unmarshalling session data")
	}

	return session, nil
}

// GetUser get the user from the token, otherwise return error
func GetUser(token *jwt.Token, db *mongo.Database) (*User, error) {
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(string)
	sessionId := claims["session_id"].(string)

	sId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return nil, err
	}

	session := Session{ID: sId}
	s, err := session.FromID(db)
	if err != nil {
		return nil, err
	}

	if time.Now().After(s.Expiry) {
		return nil, errors.New("token expired")
	}

	uId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	if uId != s.UserID {
		return nil, errors.New("user id mismatch")
	}

	user := User{ID: s.UserID}
	fromID, err := user.FromID(db)
	if err != nil {
		return nil, err
	}

	return fromID, nil
}
