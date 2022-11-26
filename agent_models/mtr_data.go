package agent_models

import (
	"time"
)

type MtrTarget struct {
	Address string `json:"address"bson:"address"`
	Result  struct {
		Triggered      bool               `json:"triggered"bson:"triggered"`
		Metrics        map[int]MtrMetrics `json:"metrics"bson:"metrics"`
		StartTimestamp time.Time          `json:"start_timestamp"bson:"start_timestamp"`
		StopTimestamp  time.Time          `json:"stop_timestamp"bson:"stop_timestamp"`
	} `json:"result"bson:"result"`
}

type MtrMetrics struct {
	Address  string `json:"address"bson:"address"`
	FQDN     string `bson:"fqdn"json:"fqdn"`
	Sent     int    `json:"sent"bson:"sent"`
	Received int    `json:"received"bson:"received"`
	Last     string `bson:"last"json:"last"`
	Avg      string `bson:"avg"json:"avg"`
	Best     string `bson:"best"json:"best"`
	Worst    string `bson:"worst"json:"worst"`
}
