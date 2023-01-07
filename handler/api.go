package handler

type ApiRequest struct {
	PIN   string      `json:"pin,omitempty"`
	ID    string      `json:"id,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
