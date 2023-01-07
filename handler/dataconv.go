package handler

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
