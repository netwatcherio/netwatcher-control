package control_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MtrData struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"json:"id"`
	Agent     primitive.ObjectID `bson:"agent"json:"agent"`
	Data      []RealMtrData      `bson:"data"json:"data"`
	Timestamp time.Time          `bson:"timestamp"json:"timestamp"`
}

type RealMtrData struct {
	Address string    `json:"address"bson:"address"`
	Result  MtrResult `json:"result"bson:"result"`
}

type MtrResult struct {
	Triggered      bool      `json:"triggered"bson:"triggered"`
	Mtr            Mtr       `json:"mtr"bson:"mtr"`
	StartTimestamp time.Time `json:"start_timestamp"bson:"start_timestamp"`
	StopTimestamp  time.Time `json:"stop_timestamp"bson:"stop_timestamp"`
}

type Mtr struct {
	Destination string        `json:"destination"bson:"destination"`
	Statistic   map[int]Stats `json:"statistic"bson:"statistic"`
}

type Stats struct {
	Sent        int     `json:"sent"bson:"sent"`
	TTL         int     `json:"ttl"bson:"TTL"`
	Target      string  `json:"target"bson:"target"`
	LastMs      float32 `json:"last_ms"bson:"last_ms"`
	BestMs      float32 `json:"best_ms"bson:"best_ms"`
	WorstMs     float32 `json:"worst_ms"bson:"worst_ms"`
	AvgMs       float32 `json:"avg_ms"bson:"avg_ms"`
	LossPercent int     `json:"loss_percent"bson:"loss_percent"`
}
