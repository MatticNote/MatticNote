package settings

import (
	"context"
	"database/sql"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"time"
)

type signinLogStruct struct {
	TriedAt   time.Time
	IsSuccess bool
	FromIp    sql.NullString
}

func securityPageGet(c *fiber.Ctx) error {
	return securityPageView(c)
}

func securityPageView(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	var signInLogs []signinLogStruct

	rows, err := database.DBPool.Query(
		context.Background(),
		"select tried_at, is_success, from_ip from signin where target_user = $1 order by tried_at desc limit 20;",
		uuid.MustParse(claim["sub"].(string)).String(),
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Render(
				"account_settings/security",
				fiber.Map{
					"SignInLogs": signInLogs,
				},
				"_layout/account",
			)
		} else {
			return err
		}
	}

	for rows.Next() {
		data := signinLogStruct{}
		err = rows.Scan(
			&data.TriedAt,
			&data.IsSuccess,
			&data.FromIp,
		)
		if err != nil {
			return err
		}
		signInLogs = append(signInLogs, data)
	}

	return c.Render(
		"account_settings/security",
		fiber.Map{
			"SignInLogs": signInLogs,
		},
		"_layout/account",
	)
}
