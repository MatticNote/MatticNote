package server

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/server/ap"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

func userProfileController(c *fiber.Ctx) error {
	if misc.IsAPAcceptHeader(c) {
		// ActivityPub Render
		targetUuid, err := internal.GetLocalUserUUIDFromUsername(c.Params("username"))
		if err != nil && err == internal.ErrNoSuchUser {
			return fiber.ErrNotFound
		} else if err != nil {
			return err
		}
		return ap.RenderUser(c, *targetUuid)
	} else {
		// Normal render
		return userProfileView(c)
	}
}

func userProfileView(c *fiber.Ctx) error {
	var (
		username       string
		displayName    misc.NullString
		summary        misc.NullString
		createdAt      misc.NullTime
		updatedAt      misc.NullTime
		isActive       bool
		isSuspend      bool
		isSuperUser    bool
		isBot          bool
		acceptManually bool
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select username, display_name, summary, created_at, updated_at, is_active, is_suspend, is_superuser, is_bot, accept_manually "+
			"from \"user\" where username ilike $1 and host is null;",
		c.Params("username"),
	).Scan(
		&username,
		&displayName,
		&summary,
		&createdAt,
		&updatedAt,
		&isActive,
		&isSuspend,
		&isSuperUser,
		&isBot,
		&acceptManually,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.ErrNotFound
		} else {
			return err
		}
	}

	if !isActive {
		return fiber.ErrGone
	}

	if isSuspend {
		return fiber.ErrForbidden
	}

	return c.Render(
		"user_profile",
		fiber.Map{
			"username":       username,
			"displayName":    displayName.String,
			"summary":        summary.String,
			"createdAt":      createdAt.Time,
			"updatedAt":      updatedAt.Time,
			"isSuperUser":    isSuperUser,
			"isBot":          isBot,
			"acceptManually": acceptManually,
		},
	)
}
