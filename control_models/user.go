package control_models

// TODO if generating db and stuff for first time,
// generate default admin and paste to console

type User struct {
	Email    string   `json:"email"`    // username
	Password string   `json:"password"` // password in sha256?
	Name     string   `json:"name"`
	Sites    []string `json:"sites"` // _id's of mongo objects
}
