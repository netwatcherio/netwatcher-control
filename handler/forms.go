package handler

type CheckNewForm struct {
	Type              string `json:"type"form:"type"`
	Target            string `json:"target"form:"target"`
	RperfServerEnable bool   `json:"target"form:"rperfServerEnable"`
	Duration          int    `json:"omitempty'"form:"duration"`
	Count             int    `json:"count,omitempty"form:"count"`
}
