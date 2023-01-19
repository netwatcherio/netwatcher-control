package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/netwatcherio/netwatcher-control/handler/agent"
	"github.com/netwatcherio/netwatcher-control/handler/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) deleteCheck() {
	r.App.Get("/delete_check/:check?", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "Application/json"
		t := c.Locals("user").(*jwt.Token)
		_, err := auth.GetUser(t, r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		aId, err := primitive.ObjectIDFromHex(c.Params("check"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// delete check data
		acc := agent.Check{ID: aId}
		err = acc.Delete(r.DB)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

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
		req := agent.CheckRequest{}
		err = c.BodyParser(&req)
		if err != nil {
			return err
		}

		check := agent.Check{ID: cId}
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
		case agent.CtMtr:
			var mtrD []*agent.MtrResult
			for _, d := range data {
				if d.Type == agent.CtMtr {
					mtr, err := d.ConvMtr()
					if err != nil {
						return err
					}
					mtrD = append(mtrD, mtr)
				}
			}
			checkData = mtrD
		case agent.CtNetworkInfo:
			var netinfoD []*agent.NetResult
			for _, d := range data {
				if d.Type == agent.CtNetworkInfo {
					netinfo, err := d.ConvNetResult()
					if err != nil {
						return err
					}
					netinfoD = append(netinfoD, netinfo)
				}
			}
			checkData = netinfoD
		case agent.CtSpeedTest:
			var speedD []*agent.SpeedTestResult
			for _, d := range data {
				if d.Type == agent.CtSpeedTest {
					speed, err := d.ConvSpeedTest()
					if err != nil {
						return err
					}
					speedD = append(speedD, speed)
				}
			}
			checkData = speedD
		case agent.CtRPerf:
			var rperfD []*agent.RPerfResults
			for _, d := range data {
				if d.Type == agent.CtRPerf {
					rperf, err := d.ConvRPerf()
					if err != nil {
						return err
					}
					rperfD = append(rperfD, rperf)
				}
			}
			checkData = rperfD
		case agent.CtPing:
			var pingD []*agent.PingResult
			for _, d := range data {
				if d.Type == agent.CtPing {
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
		req := agent.Check{}
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

		check := agent.Check{ID: cId}
		cc, err := check.Get(r.DB)
		if err != nil {
			return c.JSON(err)
		}

		return c.JSON(cc)
	})
}

// agentChecks := handler.AgentCheck{AgentID: agent.ID}
//		all, err := agentagent.GetAll(r.DB)
//		if err != nil {
//			log.Error(err)
//		}
