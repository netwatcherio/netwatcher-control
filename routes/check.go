package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/netwatcherio/netwatcher-control/handler"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html"
	"time"
)

// TODO authenticate & verify that the user is infact apart of the site etc.

func (r *Router) check() {
	r.App.Get("/check/:check?", func(c *fiber.Ctx) error {
		user, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("check") == "" {
			return c.RedirectBack("/agents")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("check"))
		if err != nil {
			return c.RedirectBack("/home")
		}

		ac := handler.AgentCheck{ID: objId}
		_, err = ac.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/agents")
		}

		agent := handler.Agent{ID: ac.AgentID}
		err = agent.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		site := handler.Site{ID: agent.Site}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		marshal, err := json.Marshal(agent)
		if err != nil {
			log.Errorf("13 %s", err)
		}

		var data []*handler.CheckData

		// todo handle time search
		if ac.Type == handler.CtRperf {
			data, err = ac.GetData(120, false, false,
				time.Now().Add(-(time.Hour * 24)), time.Now(), r.DB)
			if err != nil {
				return err
			}
		} else {
			data, err = ac.GetData(10, true, true, time.Time{}, time.Time{}, r.DB)
			if err != nil {
				return err
			}
		}

		var checkData interface{}

		switch ac.Type {
		case handler.CtMtr:
			var mtrD []*handler.MtrResult
			for _, d := range data {
				if d.Type == handler.CtMtr {
					mtr, err := d.ConvMtr()
					if err != nil {
						return err
					}
					mtrD = append(mtrD, mtr)
				}
			}
			checkData = mtrD
		case handler.CtNetinfo:
			var netinfoD []*handler.NetResult
			for _, d := range data {
				if d.Type == handler.CtNetinfo {
					netinfo, err := d.ConvNetresult()
					if err != nil {
						return err
					}
					netinfoD = append(netinfoD, netinfo)
				}
			}
			checkData = netinfoD
		case handler.CtSpeedtest:
			var speedD []*handler.SpeedTest
			for _, d := range data {
				if d.Type == handler.CtSpeedtest {
					speed, err := d.ConvSpeedtest()
					if err != nil {
						return err
					}
					speedD = append(speedD, speed)
				}
			}
			checkData = speedD
		case handler.CtRperf:
			var rperfD []*handler.RPerfResults
			for _, d := range data {
				if d.Type == handler.CtRperf {
					rperf, err := d.ConvRperf()
					if err != nil {
						return err
					}
					rperfD = append(rperfD, rperf)
				}
			}
			checkData = rperfD
		}

		bytes, err := json.Marshal(checkData)
		if err != nil {
			log.Error(err)
		}

		acBytes, err := json.Marshal(ac)
		if err != nil {
			log.Error(err)
		}

		hasData := len(data) > 0

		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("check", fiber.Map{
			"title":        agent.Name,
			"siteSelected": true,
			"siteName":     site.Name,
			"siteId":       site.ID.Hex(),
			"agents":       html.UnescapeString(string(marshal)),
			"check":        html.UnescapeString(string(acBytes)),
			"checkData":    html.UnescapeString(string(bytes)),
			"hasData":      hasData,
			"agentId":      agent.ID.Hex(),
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
		},
			"layouts/main")
	})
}

func (r *Router) checkNew() {
	r.App.Get("/check/new/:agentid?", func(c *fiber.Ctx) error {
		user, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("agentid") == "" {
			return c.Redirect("/home")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("agentid"))
		if err != nil {
			return c.Redirect("/home")
		}

		agent := handler.Agent{ID: objId}
		err = agent.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		site := handler.Site{ID: agent.Site}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		// TODO process if they are logged in or not, otherwise send them to registration/login
		return c.Render("check_new", fiber.Map{
			"title":        "new check",
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"agentId":      agent.ID.Hex(),
			"siteName":     site.Name,
			"siteSelected": true,
		},
			"layouts/main")
	})
	r.App.Post("/check/new/:agentid?", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded") // "Application/json"

		// todo recevied body is in url format, need to convert to new struct??
		//

		_, err := validateUser(r, c)
		if err != nil {
			return err
		}

		if c.Params("agentid") == "" {
			return c.Redirect("/agents")
		}
		objId, err := primitive.ObjectIDFromHex(c.Params("agentid"))
		if err != nil {
			return c.Redirect("/agents")
		}

		agent := handler.Agent{ID: objId}
		err = agent.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		site := handler.Site{ID: agent.Site}
		err = site.Get(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/home")
		}

		cCheck := new(handler.CheckNewForm)
		if err := c.BodyParser(cCheck); err != nil {
			log.Warnf("4 %s", err)
			return err
		}

		var aC handler.AgentCheck

		// todo validate target depending on check
		// ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$

		if cCheck.Type == string(handler.CtMtr) {
			if cCheck.Duration < 30 {
				cCheck.Duration = 30
			}

			aC = handler.AgentCheck{
				Type:     handler.CtMtr,
				Target:   cCheck.Target,
				AgentID:  agent.ID,
				Duration: cCheck.Duration,
			}
		} else if cCheck.Type == string(handler.CtRperf) {
			aC = handler.AgentCheck{
				Type:    handler.CtRperf,
				Target:  cCheck.Target,
				AgentID: agent.ID,
				Server:  cCheck.RperfServerEnable,
			}
		}

		err = aC.Create(r.DB)
		if err != nil {
			log.Error(err)
			return c.Redirect("/agents")
		}

		// todo create default checks such as network info and that sort of thing

		// todo handle error/success and return to home also display message for error if error
		return c.Redirect("/agent/" + agent.ID.Hex())
	})
}
