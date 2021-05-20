package account

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type issueConfirmFormStruct struct {
	Email string `validate:"required,email"`
}

func issueConfirmGet(c *fiber.Ctx) error {
	if c.Cookies(internal.JWTAuthCookieName, "") != "" {
		return c.Redirect("/web/", 307)
	}

	return issueConfirmView(c)
}

func issueConfirmView(c *fiber.Ctx, errors ...string) error {
	return c.Status(http.StatusOK).Render(
		"account/issue_confirm",
		fiber.Map{
			"CSRFFormName": csrfFormName,
			"CSRFToken":    c.Context().UserValue(csrfContextKey).(string),
			"errors":       errors,
		},
		"_layout/account",
	)
}

func issueConfirmPost(c *fiber.Ctx) error {
	formData := new(issueConfirmFormStruct)
	issueData := new(struct {
		uuidRaw    string
		email      string
		isVerified bool
	})

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return forgotPasswordView(c, errs...)
	}

	err := database.DBPool.QueryRow(
		context.Background(),
		"select \"user\".uuid, email, is_verified from \"user\" left join user_mail um on \"user\".uuid = um.uuid where email ilike $1;",
		formData.Email,
	).Scan(
		&issueData.uuidRaw,
		&issueData.email,
		&issueData.isVerified,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Redirect("/account/login?issued_confirm=true")
		} else {
			return err
		}
	}

	if !issueData.isVerified {
		if err := internal.IssueVerifyEmail(uuid.MustParse(issueData.uuidRaw), issueData.email); err != nil {
			return err
		}
	}

	return c.Redirect("/account/login?issued_confirm=true")
}
