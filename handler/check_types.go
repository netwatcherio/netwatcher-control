package handler

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type CheckType string

const (
	CtRperf     CheckType = "RPERF"
	CtMtr       CheckType = "MTR"
	CtSpeedtest CheckType = "SPEEDTEST"
	CtNetinfo   CheckType = "NETINFO"
)

func (cd *CheckData) ConvNetresult() (*NetResult, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r NetResult
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *CheckData) ConvMtr() (*MtrResult, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r MtrResult
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *CheckData) ConvSpeedtest() (*SpeedTest, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r SpeedTest
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *CheckData) ConvRperf() (*RPerfResults, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r RPerfResults
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

type MtrResult struct {
	StartTimestamp time.Time `json:"start_timestamp"bson:"start_timestamp"`
	StopTimestamp  time.Time `json:"stop_timestamp"bson:"stop_timestamp"`
	Triggered      bool      `bson:"triggered"json:"triggered"`
	Report         struct {
		Mtr struct {
			Src        string `json:"src"bson:"src"`
			Dst        string `json:"dst"bson:"dst"`
			Tos        int    `json:"tos"bson:"tos"`
			Tests      int    `json:"tests"bson:"tests"`
			Psize      string `json:"psize"bson:"psize"`
			Bitpattern string `json:"bitpattern"bson:"bitpattern"`
		} `json:"mtr"bson:"mtr"`
		Hubs []struct {
			Count int     `json:"count"bson:"count"`
			Host  string  `json:"host"bson:"host"`
			ASN   string  `json:"ASN"bson:"ASN"`
			Loss  float64 `json:"Loss%"bson:"Loss%"`
			Drop  int     `json:"Drop"bson:"Drop"`
			Rcv   int     `json:"Rcv"bson:"Rcv"`
			Snt   int     `json:"Snt"bson:"Snt"`
			Best  float64 `json:"Best"bson:"Best"`
			Avg   float64 `json:"Avg"bson:"Avg"`
			Wrst  float64 `json:"Wrst"bson:"Wrst"`
			StDev float64 `json:"StDev"bson:"StDev"`
			Gmean float64 `json:"Gmean"bson:"Gmean"`
			Jttr  float64 `json:"Jttr"bson:"Jttr"`
			Javg  float64 `json:"Javg"bson:"Javg"`
			Jmax  float64 `json:"Jmax"bson:"Jmax"`
			Jint  float64 `json:"Jint"bson:"Jint"`
		} `json:"hubs"bson:"hubs"`
	} `json:"report"bson:"report"`
}
type NetResult struct {
	LocalAddress     string    `json:"local_address"bson:"local_address"`
	DefaultGateway   string    `json:"default_gateway"bson:"default_gateway"`
	PublicAddress    string    `json:"public_address"bson:"public_address"`
	InternetProvider string    `json:"internet_provider"bson:"internet_provider"`
	Lat              string    `json:"lat"bson:"lat"`
	Long             string    `json:"long"bson:"long"`
	Timestamp        time.Time `json:"timestamp"bson:"timestamp"`
}
type SpeedTest struct {
	Latency   time.Duration `json:"latency"bson:"latency"`
	DLSpeed   float64       `json:"dl_speed"bson:"dl_speed"`
	ULSpeed   float64       `json:"ul_speed"bson:"ul_speed"`
	Server    string        `json:"server"bson:"server"`
	Host      string        `json:"host"bson:"host"`
	Timestamp time.Time     `json:"timestamp"bson:"timestamp"`
}
type RPerfResults struct {
	StartTimestamp time.Time `json:"start_timestamp"bson:"start_timestamp"`
	StopTimestamp  time.Time `json:"stop_timestamp"bson:"stop_timestamp"`
	Config         struct {
		Additional struct {
			IpVersion   int  `json:"ip_version"bson:"ip_version"`
			OmitSeconds int  `json:"omit_seconds"bson:"omit_seconds"`
			Reverse     bool `json:"reverse"bson:"reverse"`
		} `json:"additional"bson:"additional"`
		Common struct {
			Family  string `json:"family"bson:"family"`
			Length  int    `json:"length"bson:"length"`
			Streams int    `json:"streams"bson:"streams"`
		} `json:"common"bson:"common"`
		Download struct {
		} `json:"download"bson:"download"`
		Upload struct {
			Bandwidth    int     `json:"bandwidth"bson:"bandwidth"`
			Duration     float64 `json:"duration"bson:"duration"`
			SendInterval float64 `json:"send_interval"bson:"send_interval"`
		} `json:"upload"bson:"upload"`
	} `json:"config"bson:"config"`
	Streams []struct {
		Abandoned bool `json:"abandoned"bson:"abandoned"`
		Failed    bool `json:"failed"bson:"failed"`
		Intervals struct {
			Receive []struct {
				BytesReceived     int     `json:"bytes_received"bson:"bytes_received"`
				Duration          float64 `json:"duration"bson:"duration"`
				JitterSeconds     float64 `json:"jitter_seconds"bson:"jitter_seconds"`
				PacketsDuplicated int     `json:"packets_duplicated"bson:"packets_duplicated"`
				PacketsLost       int     `json:"packets_lost"bson:"packets_lost"`
				PacketsOutOfOrder int     `json:"packets_out_of_order"bson:"packets_out_of_order"`
				PacketsReceived   int     `json:"packets_received"bson:"packets_received"`
				Timestamp         float64 `json:"timestamp"bson:"timestamp"`
				UnbrokenSequence  int     `json:"unbroken_sequence"bson:"unbroken_sequence"`
			} `json:"receive"bson:"receive"`
			Send []struct {
				BytesSent    int     `json:"bytes_sent"bson:"bytes_sent"`
				Duration     float64 `json:"duration"bson:"duration"`
				PacketsSent  int     `json:"packets_sent"bson:"packets_sent"`
				SendsBlocked int     `json:"sends_blocked"bson:"sends_blocked"`
				Timestamp    float64 `json:"timestamp"bson:"timestamp"`
			} `json:"send"bson:"send"`
			Summary struct {
				BytesReceived            int     `json:"bytes_received"bson:"bytes_received"`
				BytesSent                int     `json:"bytes_sent"bson:"bytes_sent"`
				DurationReceive          float64 `json:"duration_receive"bson:"duration_receive"`
				DurationSend             float64 `json:"duration_send"bson:"duration_send"`
				FramedPacketSize         int     `json:"framed_packet_size"bson:"framed_packet_size"`
				JitterAverage            float64 `json:"jitter_average"bson:"jitter_average"`
				JitterPacketsConsecutive int     `json:"jitter_packets_consecutive"bson:"jitter_packets_consecutive"`
				PacketsDuplicated        int     `json:"packets_duplicated"bson:"packets_duplicated"`
				PacketsLost              int     `json:"packets_lost"bson:"packets_lost"`
				PacketsOutOfOrder        int     `json:"packets_out_of_order"bson:"packets_out_of_order"`
				PacketsReceived          int     `json:"packets_received"bson:"packets_received"`
				PacketsSent              int     `json:"packets_sent"bson:"packets_sent"`
			} `json:"summary"bson:"summary"`
		} `json:"intervals"bson:"intervals"`
	} `json:"streams"bson:"streams"`
	Success bool `json:"success"bson:"success"`
	Summary struct {
		BytesReceived            int     `json:"bytes_received"bson:"bytes_received"`
		BytesSent                int     `json:"bytes_sent"bson:"bytes_sent"`
		DurationReceive          float64 `json:"duration_receive"bson:"duration_receive"`
		DurationSend             float64 `json:"duration_send"bson:"duration_send"`
		FramedPacketSize         int     `json:"framed_packet_size"bson:"framed_packet_size"`
		JitterAverage            float64 `json:"jitter_average"bson:"jitter_average"`
		JitterPacketsConsecutive int     `json:"jitter_packets_consecutive"bson:"jitter_packets_consecutive"`
		PacketsDuplicated        int     `json:"packets_duplicated"bson:"packets_duplicated"`
		PacketsLost              int     `json:"packets_lost"bson:"packets_lost"`
		PacketsOutOfOrder        int     `json:"packets_out_of_order"bson:"packets_out_of_order"`
		PacketsReceived          int     `json:"packets_received"bson:"packets_received"`
		PacketsSent              int     `json:"packets_sent"bson:"packets_sent"`
	} `json:"summary"bson:"summary"`
}
