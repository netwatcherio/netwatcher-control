package handler

type Data struct {
	PIN    string      `json:"pin,omitempty"`
	ID     string      `json:"id,omitempty"`
	Checks interface{} `json:"checks"`
	Error  string      `json:"error,omitempty"`
}
