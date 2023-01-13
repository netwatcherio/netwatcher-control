package checks

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (cd *Data) ConvNetResult() (*NetResult, error) {
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

func (cd *Data) ConvMtr() (*MtrResult, error) {
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

func (cd *Data) ConvPing() (*PingResult, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r PingResult
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *Data) ConvSpeedTest() (*SpeedTestResult, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r SpeedTestResult
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *Data) ConvRPerf() (*RPerfResults, error) {
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
