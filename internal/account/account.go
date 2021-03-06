package account

import (
	"database/sql"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserGone           = errors.New("user is gone")
)

const (
	TokenCookieName = "mn_token"
)

func AuthenticateUser(
	email,
	password string,
) (*schemas.User, error) {
	var (
		userId       ksuid.KSUID
		userPassword []byte
	)

	err := database.Database.QueryRow(
		"SELECT u.id, ua.password FROM users_email "+
			"LEFT OUTER JOIN users u on u.id = users_email.id "+
			"LEFT JOIN users_auth ua on u.id = ua.id "+
			"WHERE email = $1 AND (deleted_at IS NULL OR deleted_at > now());",
		email,
	).
		Scan(&userId, &userPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
	}

	err = bcrypt.CompareHashAndPassword(userPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid && user.DeletedAt.Time.Before(time.Now()) {
		return nil, ErrUserGone
	}

	return user, nil
}

func InsertTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		Secure:   false,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func GetUser(userId ksuid.KSUID) (*schemas.User, error) {
	user := new(schemas.User)

	err := database.Database.QueryRow(
		"SELECT id, username, host, display_name, headline, description, created_at, is_silence, is_suspend, is_moderator, is_admin, deleted_at FROM users WHERE id = $1;", userId.String()).
		Scan(&user.ID, &user.Username, &user.Host, &user.DisplayName, &user.Headline, &user.Description, &user.CreatedAt, &user.IsSilence, &user.IsSuspend, &user.IsModerator, &user.IsAdmin, &user.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		} else {
			return nil, err
		}
	}

	return user, nil
}

func GetUserByUsername(username string, host ...string) (*schemas.User, error) {
	var (
		userId ksuid.KSUID
		err    error
	)

	if len(host) > 0 {
		err = database.Database.QueryRow("SELECT id FROM users WHERE username ILIKE $1 AND host ILIKE $2", username, host[0]).Scan(&userId)
	} else {
		err = database.Database.QueryRow("SELECT id FROM users WHERE username ILIKE $1 AND host IS NULL", username).Scan(&userId)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		} else {
			return nil, err
		}
	}

	return GetUser(userId)
}

func GetUserEmail(userId ksuid.KSUID) (string, error) {
	var email string
	err := database.Database.QueryRow("SELECT email FROM	users_email WHERE id = $1", userId.String()).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		} else {
			return "", err
		}
	}

	return email, nil
}
