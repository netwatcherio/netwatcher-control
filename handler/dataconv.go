package handler

import (
	"github.com/netwatcherio/netwatcher-agent/checks"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (cd *CheckData) ConvNetresult() (*checks.NetResult, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r checks.NetResult
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *CheckData) ConvMtr() (*checks.MtrResult, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r checks.MtrResult
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *CheckData) ConvSpeedtest() (*checks.SpeedTest, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r checks.SpeedTest
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}

func (cd *CheckData) ConvRperf() (*checks.RPerfResults, error) {
	crM, err := bson.Marshal(cd.Result)
	if err != nil {
		log.Error(err)
	}

	var r checks.RPerfResults
	err = bson.Unmarshal(crM, &r)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &r, nil
}
