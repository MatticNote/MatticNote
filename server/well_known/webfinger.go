package well_known

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"regexp"
)

func webfinger(c *fiber.Ctx) error {
	resource := c.Query("resource")
	if resource == "" {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	rRegex := regexp.MustCompile(`(acct:)?([a-zA-Z0-9_]+)@(.+)`)
	rResult := rRegex.FindStringSubmatch(resource)

	if len(rResult) != 4 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	endpointUrl, err := url.Parse(config.Config.Server.Endpoint)
	if err != nil {
		return err
	}
	if endpointUrl.Host != rResult[3] {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	targetUser, err := user.GetLocalUserFromUsername(rResult[2])
	if err != nil {
		switch err {
		case user.ErrNoSuchUser:
			c.Status(fiber.StatusNotFound)
		case user.ErrUserGone:
			c.Status(fiber.StatusGone)
		default:
			return err
		}
		return nil
	}

	wfJsonMarshal, err := json.Marshal(fiber.Map{
		"subject": fmt.Sprintf("acct:%s@%s", targetUser.Username, endpointUrl.Host),
		"aliases": []interface{}{
			fmt.Sprintf("%s/@%s", config.Config.Server.Endpoint, targetUser.Username),
			fmt.Sprintf("%s/user/%s", config.Config.Server.Endpoint, targetUser.Username),
		},
		"links": []fiber.Map{
			{
				"rel":  "http://webfinger.net/rel/profile-page",
				"type": "text/html",
				"href": fmt.Sprintf("%s/@%s", config.Config.Server.Endpoint, targetUser.Username),
			},
			{
				"rel":  "self",
				"type": "application/activity+json",
				"href": fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetUser.Uuid),
			},
		},
	})

	c.Set("Content-Type", "application/jrd+json")
	return c.Send(wfJsonMarshal)
}
