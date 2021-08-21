package ap

import (
	"github.com/MatticNote/MatticNote/worker"
	"github.com/gocraft/work"
	"github.com/gofiber/fiber/v2"
)

func inboxPost(c *fiber.Ctx) error {
	_, err := worker.Enqueue.Enqueue("inbox_worker", work.Q{})
	if err != nil {
		return err
	}
	c.Status(fiber.StatusAccepted)
	return nil
}
