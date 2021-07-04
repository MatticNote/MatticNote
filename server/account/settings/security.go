package settings

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"image/png"
	"time"
)

type signinLogStruct struct {
	TriedAt   time.Time
	IsSuccess bool
	FromIp    sql.NullString
}

type user2faSetupForm struct {
	Token string `validate:"required,len=6,containsany=0123456789"`
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
				"_layout/settings",
			)
		} else {
			return err
		}
	}
	defer rows.Close()

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

	if rows.Err() != nil {
		return rows.Err()
	}

	return c.Render(
		"account_settings/security",
		fiber.Map{
			"SignInLogs": signInLogs,
		},
		"_layout/settings",
	)
}

func setup2faGet(c *fiber.Ctx) error {
	return setup2faView(c, false)
}

func setup2faView(c *fiber.Ctx, isFail bool) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	totpCode, err := internal.Setup2faCode(uuid.MustParse(claim["sub"].(string)))
	if err != nil {
		if err == internal.Err2faAlreadyEnabled {
			return c.Redirect("/account/settings/security")
		} else {
			return err
		}
	}

	image, err := totpCode.Image(256, 256)
	if err != nil {
		return err
	}
	var qrBytes bytes.Buffer
	err = png.Encode(&qrBytes, image)
	if err != nil {
		return err
	}

	return c.Render(
		"account_settings/security_2fa_setup",
		fiber.Map{
			"QRCode":       fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrBytes.Bytes())),
			"SecretCode":   totpCode.Secret(),
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
			"IsFail":       isFail,
		},
		"_layout/settings",
	)
}

func setup2faPost(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	formData := new(user2faSetupForm)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(formData); errs != nil {
		return setup2faView(c, true)
	}

	targetUuid := uuid.MustParse(claim["sub"].(string))

	err := internal.Validate2faCode(targetUuid, formData.Token)
	if err != nil {
		if err == internal.ErrInvalid2faToken {
			return setup2faView(c, true)
		} else {
			return err
		}
	}

	if err := internal.Enable2faAuth(targetUuid); err != nil {
		return err
	}

	return c.Redirect("/account/settings/security/2fa/backup")
}

func get2faBackup(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	code, err := internal.Get2faBackupCode(uuid.MustParse(claim["sub"].(string)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Redirect("/account/settings/security/2fa/setup", fiber.StatusTemporaryRedirect)
		} else {
			return err
		}
	}

	return c.Render(
		"account_settings/security_2fa_backup",
		fiber.Map{
			"BackupCodes":  code,
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
		},
		"_layout/settings",
	)
}

func regenerate2faBackup(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	err := internal.Regenerate2faBackupCode(uuid.MustParse(claim["sub"].(string)))
	if err != nil {
		return err
	}

	return get2faBackup(c)
}
