package well_known

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"regexp"
)

func webfinger(c *fiber.Ctx) error {
	resource := c.Query("resource")
	if resource == "" {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	rRegex := regexp.MustCompile(`acct:([a-zA-Z0-9_]+)@(.+)`)
	rResult := rRegex.FindStringSubmatch(resource)

	targetUser, err := internal.GetLocalUserFromUsername(rResult[1])
	if err != nil {
		switch err {
		case internal.ErrNoSuchUser:
			c.Status(fiber.StatusNotFound)
		case internal.ErrUserGone:
			c.Status(fiber.StatusGone)
		case internal.ErrUserSuspended:
			c.Status(fiber.StatusForbidden)
		default:
			return err
		}
		return nil
	}

	wfJsonMarshal, err := json.Marshal(fiber.Map{
		"subject": fmt.Sprintf("acct:%s@%s", targetUser.Username, rResult[2]),
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
