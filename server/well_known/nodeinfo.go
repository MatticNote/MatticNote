package well_known

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/gofiber/fiber/v2"
)

func nodeinfoWK(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"links": []fiber.Map{
			{
				"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.1",
				"href": fmt.Sprintf("%s/nodeinfo/2.1", config.Config.Server.Endpoint),
			},
		},
	})
}
