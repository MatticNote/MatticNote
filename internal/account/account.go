package account

import (
	"database/sql"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

const (
	TokenCookieName = "mn_token"
)

func AuthenticateUser(
	email,
	password string,
) (*types.User, error) {
	var (
		user         types.User
		userPassword sql.NullString
	)

	err := database.Database.QueryRow(
		"SELECT "+
			"\"user\".id, "+
			"username, "+
			"display_name, "+
			"headline, "+
			"description, "+
			"created_at, "+
			"is_silence, "+
			"is_suspend, "+
			"is_active, "+
			"is_moderator, "+
			"is_admin, "+
			"email, "+
			"verified, "+
			"password "+
			"FROM \"user\" "+
			"LEFT JOIN user_email ue "+
			"ON \"user\".id = ue.id "+
			"LEFT JOIN user_auth ua ON "+
			"\"user\".id = ua.id "+
			"WHERE email ILIKE $1 AND "+
			"host IS NULL;",
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Headline,
		&user.Description,
		&user.CreatedAt,
		&user.IsSilence,
		&user.IsSuspend,
		&user.IsActive,
		&user.IsModerator,
		&user.IsAdmin,
		&user.Email,
		&user.EmailVerified,
		&userPassword,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	if !userPassword.Valid {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword.String), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func InsertTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		Secure:   false,
		HTTPOnly: true,
		SameSite: "Strict",
	})
}

func GetUser(userId ksuid.KSUID) (*types.User, error) {
	var user types.User

	err := database.Database.QueryRow(
		"SELECT "+
			"\"user\".id, "+
			"username, "+
			"display_name, "+
			"headline, "+
			"description, "+
			"created_at, "+
			"is_silence, "+
			"is_suspend, "+
			"is_active, "+
			"is_moderator, "+
			"is_admin, "+
			"email, "+
			"verified "+
			"FROM \"user\" "+
			"LEFT JOIN user_email ue "+
			"ON \"user\".id = ue.id "+
			"WHERE \"user\".id = $1",
		userId.String(),
	).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Headline,
		&user.Description,
		&user.CreatedAt,
		&user.IsSilence,
		&user.IsSuspend,
		&user.IsActive,
		&user.IsModerator,
		&user.IsAdmin,
		&user.Email,
		&user.EmailVerified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		} else {
			return nil, err
		}
	}

	return &user, nil
}
