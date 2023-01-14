package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/netwatcherio/netwatcher-control/handler/agent/checks"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) getCheckData() {
	r.App.Post("/check/:check?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		cId, err := primitive.ObjectIDFromHex(c.Params("check"))
		if err != nil {
			return c.JSON(err)
		}

		// require check request
		req := checks.CheckRequest{}
		err = c.BodyParser(&req)
		if err != nil {
			return err
		}

		check := checks.Check{ID: cId}
		_, err = check.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		data, err := check.GetData(req, r.DB)
		if err != nil {
			log.Errorf("no data found")
		}

		var checkData interface{}

		switch check.Type {
		case checks.CtMtr:
			var mtrD []*checks.MtrResult
			for _, d := range data {
				if d.Type == checks.CtMtr {
					mtr, err := d.ConvMtr()
					if err != nil {
						return err
					}
					mtrD = append(mtrD, mtr)
				}
			}
			checkData = mtrD
		case checks.CtNetworkInfo:
			var netinfoD []*checks.NetResult
			for _, d := range data {
				if d.Type == checks.CtNetworkInfo {
					netinfo, err := d.ConvNetResult()
					if err != nil {
						return err
					}
					netinfoD = append(netinfoD, netinfo)
				}
			}
			checkData = netinfoD
		case checks.CtSpeedTest:
			var speedD []*checks.SpeedTestResult
			for _, d := range data {
				if d.Type == checks.CtSpeedTest {
					speed, err := d.ConvSpeedTest()
					if err != nil {
						return err
					}
					speedD = append(speedD, speed)
				}
			}
			checkData = speedD
		case checks.CtRPerf:
			var rperfD []*checks.RPerfResults
			for _, d := range data {
				if d.Type == checks.CtRPerf {
					rperf, err := d.ConvRPerf()
					if err != nil {
						return err
					}
					rperfD = append(rperfD, rperf)
				}
			}
			checkData = rperfD
		case checks.CtPing:
			var pingD []*checks.PingResult
			for _, d := range data {
				if d.Type == checks.CtPing {
					ping, err := d.ConvPing()
					if err != nil {
						return err
					}
					pingD = append(pingD, ping)
				}
			}
			checkData = pingD
		}

		return c.JSON(checkData)
	})
}

func (r *Router) checkNew() {
	r.App.Post("/check/new/:agentid?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("agentid"))
		if err != nil {
			return c.JSON(err)
		}

		// require check request
		req := checks.Check{}
		err = c.BodyParser(&req)
		if err != nil {
			return c.JSON(err)
		}
		req.AgentID = aId

		// todo handle edge cases? the user *could* break their install if not... hmmm...

		err = req.Create(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

func (r *Router) getCheck() {
	r.App.Get("/check/:check?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.JSON(err)
		}

		cId, err := primitive.ObjectIDFromHex(c.Params("check"))
		if err != nil {
			return c.JSON(err)
		}

		// todo handle edge cases? the user *could* break their install if not... hmmm...

		check := checks.Check{ID: cId}
		cc, err := check.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		return c.JSON(cc)
	})
}

// agentChecks := handler.AgentCheck{AgentID: agent.ID}
//		all, err := agentChecks.GetAll(r.DB)
//		if err != nil {
//			log.Error(err)
//		}
