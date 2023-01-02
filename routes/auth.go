package routes

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"netwatcher-control/handler"
	"netwatcher-control/models"
)

func (r *Router) authRegister() {
	r.App.Get("/auth/register", func(c *fiber.Ctx) error {
		b, _ := handler.ValidateSession(c, r.Session, r.DB)
		if b {
			return c.Redirect("/home")
		}
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("auth", fiber.Map{
			"title": "auth", "login": false})
	})
	r.App.Post("/auth/register", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		registerUser := new(models.RegisterUser)
		if err := c.BodyParser(registerUser); err != nil {
			log.Warnf("%s", err)
			return err
		}

		if registerUser.Password != registerUser.PasswordConfirm {
			//todo handle error and show on auth page using sessions??
			return c.Redirect("/auth/register")
		}

		pwd, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 15)
		if err != nil {
			log.Errorf("%s", err)
			return c.Redirect("/auth/login")
		}

		user := models.User{
			ID:        primitive.NewObjectID(),
			Email:     registerUser.Email,
			FirstName: registerUser.FirstName,
			LastName:  registerUser.LastName,
			Admin:     false,
			Password:  string(pwd),
			Sites:     nil,
			Verified:  false,
		}

		ucb, err2 := user.Create(r.DB)
		if err2 != nil || !ucb {
			log.Infof("%s", "error creating user")
			return c.Redirect("/auth/register")
		}

		//todo handle success and send to login page
		return c.Redirect("/auth/login")
	})
}

func (r *Router) authLogin() {
	r.App.Get("/auth/login", func(c *fiber.Ctx) error {
		b, _ := handler.ValidateSession(c, r.Session, r.DB)
		if b {
			return c.Redirect("/home")
		}
		// Render index within layouts/main
		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("auth", fiber.Map{
			"title": "auth", "login": true})
	})
	r.App.Post("/auth/login", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		loginUser := new(models.LoginUser)
		if err := c.BodyParser(loginUser); err != nil {
			log.Warnf("4 %s", err)
			return err
		}

		user := models.User{Email: loginUser.Email}

		// get user from email
		usr, err2 := user.GetUserFromEmail(r.DB)
		if err2 != nil {
			log.Warnf("3 %s", err2)
			return c.Redirect("/auth/login")
		}

		err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(loginUser.Password))
		if err != nil {
			log.Errorf("%s", err)
			return c.Redirect("/auth/login")
		}

		// create token session
		b, err := handler.LoginSession(c, r.Session, r.DB, usr.ID)
		if err != nil || !b {
			log.Warnf("5 %s, 2 %b", err, b)
			return c.Redirect("/auth/login")
		}
		// todo handle success and return to home
		return c.Redirect("/home")
	})
}

func (r *Router) authLogout() {
	r.App.Get("/logout", func(c *fiber.Ctx) error {
		handler.LogoutSession(c, r.Session)

		return c.Redirect("/auth")
	})
}

func (r *Router) auth() {
	r.App.Get("/auth", func(c *fiber.Ctx) error {
		b, _ := handler.ValidateSession(c, r.Session, r.DB)
		if b {
			return c.Redirect("/home")
		}
		return c.Redirect("/auth/login")
	})
}
