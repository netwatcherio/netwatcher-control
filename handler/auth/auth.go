package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login returns error on fail, nil on success
func (r *Login) Login(db *mongo.Database) (string, error) {
	if r.Email == "" {
		// todo validate email
		return "", errors.New("invalid email address")
	}

	u := User{Email: r.Email}
	user, err := u.FromEmail(db)
	if err != nil {
		return "", err
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(r.Password), 15)
	if err != nil {
		return "", err
	}

	if string(pwd) == user.Password {
		return "", errors.New("invalid password, please ensure passwords match")
	}

	session := Session{
		UserID: user.ID,
	}
	err = session.Create(db)
	if err != nil {
		return "", err
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id":    session.UserID.Hex(),
		"session_id": session.ID.Hex(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ping pong")) // TODO lol
	if err != nil {
		return "", err
	}

	return t, nil
}

type Register struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Company   string `json:"company"`
}

// Register returns error on fail, nil on success
func (r *Register) Register(db *mongo.Database) (string, error) {
	if r.FirstName == "" {
		return "", errors.New("invalid first name")
	}
	if r.LastName == "" {
		return "", errors.New("invalid first name")
	}
	if r.Email == "" {
		// todo validate email
		return "", errors.New("invalid email name")
	}
	if r.Password == "" {
		return "", errors.New("invalid password, please ensure passwords match")
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(r.Password), 15)
	if err != nil {
		return "", err
	}

	user := User{
		Email:           r.Email,
		FirstName:       r.FirstName,
		LastName:        r.LastName,
		Company:         r.Company,
		Admin:           false,
		Password:        string(pwd),
		Sites:           nil,
		Verified:        false,
		CreateTimestamp: time.Now(),
	}

	err = user.Create(db)
	if err != nil {
		return "", err
	}

	session := Session{
		UserID: user.ID,
	}
	err = session.Create(db)
	if err != nil {
		return "", err
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id":    session.UserID.Hex(),
		"session_id": session.ID.Hex(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("ping pong")) // TODO lol
	if err != nil {
		return "", err
	}

	return t, nil
}
