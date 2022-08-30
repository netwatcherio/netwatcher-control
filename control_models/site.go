package control_models

type Site struct {
	Agents  []string `json:"agents"` // _id of mongo objects
	Members []struct {
		User string `json:"user"` // _id
		Role int    `json:"role"`
		// roles: 0=READ ONLY, 1=READ-WRITE (Create only), 2=ADMIN (Delete Agents), 3=OWNER (Delete Sites)
		// ADMINS can regenerate agent pins
	}
}
