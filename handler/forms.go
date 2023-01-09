package handler

type CheckNewForm struct {
	Type              string `json:"type"form:"type"`
	Target            string `json:"target"form:"target"`
	RperfServerEnable bool   `json:"server"form:"rperfServerEnable"`
	Duration          int    `json:"duration'"form:"duration"`
	Count             int    `json:"count,omitempty"form:"count"`
	Interval          int    `json:"interval,omitempty"form:"interval"`
}
