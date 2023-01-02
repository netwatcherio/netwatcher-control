package handler

import (
	"encoding/json"
	"github.com/netwatcherio/netwatcher-agent/checks"
	"go.mongodb.org/mongo-driver/mongo"
	"netwatcher-control/models"
)

/*"agents":           html.UnescapeString(string(marshal)),
"mtr":              html.UnescapeString(string(marshalMtr)),
"speed":            html.UnescapeString(string(marshalSpeed)),
"online":           stats.Online,
"agentStats":       html.UnescapeString(string(doc)),
"publicAddress":    networkData.Data.PublicAddress,
"localAddress":     networkData.Data.LocalAddress,
"defaultGateway":   networkData.Data.DefaultGateway,
"internetProvider": networkData.Data.InternetProvider,
"uploadSpeed":      math.Round(speedData[0].Data.ULSpeed),
"downloadSpeed":    math.Round(speedData[0].Data.DLSpeed),
"speedtestPending": agent.AgentConfig.SpeedTestPending,*/

type AgentStats struct {
	NetInfo          checks.NetResult
	SpeedTestInfo    checks.SpeedTest
	SpeedTestPending bool
}

func GetAgentStats(agent models.Agent, db *mongo.Database) (*AgentStats, error) {
	var stats *AgentStats

	// get the latest net stats
	agentCheck := models.AgentCheck{AgentID: agent.ID, Type: models.CT_NetInfo}
	get, err := agentCheck.GetData(1, true, nil, nil, db)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(get[0])
	if err != nil {
		return nil, err
	}
	var netInfo checks.NetResult

	err = json.Unmarshal(bytes, &netInfo)
	if err != nil {
		return nil, err
	}

	stats.NetInfo = netInfo

	// todo check the agent check it self to see if the speedtest is pending, else check and add the speedtest stats

	// get the latest net stats
	agentCheck = models.AgentCheck{AgentID: agent.ID, Type: models.CT_SpeedTest}
	get, err = agentCheck.GetData(1, true, nil, nil, db)
	if err != nil {
		return nil, err
	}

	bytes, err = json.Marshal(get[0])
	if err != nil {
		return nil, err
	}
	var speedTest checks.SpeedTest

	err = json.Unmarshal(bytes, &speedTest)
	if err != nil {
		return nil, err
	}
	stats.SpeedTestInfo = speedTest

	return stats, nil
}
