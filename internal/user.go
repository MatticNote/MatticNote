package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrUserExists        = errors.New("username or email is already taken")
	ErrLoginFailed       = errors.New("there is an error in the login name or password")
	ErrEmailAuthRequired = errors.New("target user required email authentication")
	Err2faRequired       = errors.New("target user required two factor authentication") // WIP
	ErrNoSuchUser        = errors.New("target user was not found")
	ErrUserSuspended     = errors.New("target user is suspended")
	ErrInvalidPassword   = errors.New("invalid password")
)

type LocalUserStruct struct {
	Uuid           uuid.UUID
	Username       string
	Email          string
	DisplayName    misc.NullString
	Summary        misc.NullString
	CreatedAt      misc.NullTime
	UpdatedAt      misc.NullTime
	IsSilence      bool
	AcceptManually bool
	IsSuperuser    bool
	IsBot          bool
}

const (
	UserPasswordHashCost = 12
	UserKeyPairLength    = 2048
)

func RegisterLocalUser(email, username, password string, skipEmailVerify bool) (uuid.UUID, error) {
	var count int
	err := database.DBPool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM \"user\" LEFT JOIN user_mail um on \"user\".uuid = um.uuid WHERE (username ILIKE $1 AND host IS NULL) OR email ILIKE $2;",
		username,
		email,
	).Scan(&count)
	if err != nil && err != pgx.ErrNoRows {
		return uuid.Nil, err
	}

	if count > 0 {
		return uuid.Nil, ErrUserExists
	}

	newUuid := uuid.Must(uuid.NewRandom())

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return uuid.Nil, err
	}
	defer func(tx pgx.Tx) {
		_ = tx.Rollback(context.Background())
	}(tx)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO \"user\"(uuid, username) VALUES ($1, $2);",
		newUuid,
		username,
	)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_mail(uuid, email, is_verified) VALUES ($1, $2, $3);",
		newUuid,
		email,
		skipEmailVerify,
	)
	if err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), UserPasswordHashCost)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_password(uuid, password) VALUES ($1, $2);",
		newUuid,
		hashedPassword,
	)
	if err != nil {
		return uuid.Nil, err
	}

	rsaPrivateKey, rsaPublicKey := misc.GenerateRSAKeypair(UserKeyPairLength)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_signature_key(uuid, public_key, private_key) VALUES ($1, $2, $3);",
		newUuid,
		string(rsaPublicKey),
		string(rsaPrivateKey),
	)
	if err != nil {
		return uuid.Nil, err
	}

	if !skipEmailVerify {
		if err := IssueVerifyEmail(newUuid, email); err != nil {
			return uuid.Nil, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	return newUuid, err
}

func ValidateLoginUser(login, password string) (uuid.UUID, error) {
	var (
		targetUuid     uuid.UUID
		isMailVerified bool
		targetPassword []byte
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"SELECT \"user\".uuid, is_verified, password FROM \"user\", user_mail um, user_password up WHERE \"user\".uuid = um.uuid AND \"user\".uuid = up.uuid AND email ILIKE $1 OR (username ILIKE $1 AND host IS NULL);",
		login,
	).Scan(
		&targetUuid,
		&isMailVerified,
		&targetPassword,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, ErrLoginFailed
		} else {
			return uuid.Nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(targetPassword, []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return targetUuid, ErrLoginFailed
		} else {
			return uuid.Nil, err
		}
	}

	if !isMailVerified {
		return targetUuid, ErrEmailAuthRequired
	}

	return targetUuid, nil
}

func ValidateUserPassword(targetUuid uuid.UUID, password string) error {
	var (
		targetPassword []byte
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select password from user_password where uuid = $1;",
		targetUuid.String(),
	).Scan(
		&targetPassword,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNoSuchUser
		} else {
			return err
		}
	}

	if err := bcrypt.CompareHashAndPassword(targetPassword, []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrInvalidPassword
		} else {
			return err
		}
	}

	return nil
}

func GetLocalUser(targetUuid string) (*LocalUserStruct, error) {
	targetUuidParsed, err := uuid.Parse(targetUuid)
	if err != nil {
		return nil, err
	}

	target := new(LocalUserStruct)
	var isSuspend bool

	err = database.DBPool.QueryRow(
		context.Background(),
		"select \"user\".uuid, username, email, display_name, summary, created_at, updated_at, is_silence, accept_manually, is_superuser, is_bot, is_suspend from \"user\" left join user_mail um on \"user\".uuid = um.uuid where \"user\".uuid = $1",
		targetUuidParsed.String(),
	).Scan(
		&target.Uuid,
		&target.Username,
		&target.Email,
		&target.DisplayName,
		&target.Summary,
		&target.CreatedAt,
		&target.UpdatedAt,
		&target.IsSilence,
		&target.AcceptManually,
		&target.IsSuperuser,
		&target.IsBot,
		&isSuspend,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoSuchUser
		} else {
			return nil, err
		}
	}

	if isSuspend {
		return nil, ErrUserSuspended
	}

	return target, nil
}

func IssueForgotPassword(targetEmail string) error {
	verifyToken := misc.GenToken(64)

	cmdTag, err := database.DBPool.Exec(
		context.Background(),
		"insert into user_reset_password(key, target) values ($1, (select user_mail.uuid from user_mail where is_verified is true and email ilike $2)) "+
			"on conflict on constraint user_reset_password_target do update set key = $1, expired = default;",
		verifyToken,
		targetEmail,
	)
	if err == nil {
		if cmdTag.RowsAffected() > 0 {
			// TODO: 将来的に"template/text"を使う
			body := fmt.Sprintf(
				"Hello. This is MatticNote.\n\n"+
					"We have received a request to reset the password of the account registered with this email address.\n\n"+
					"You can reset your password by clicking the URL below: \n\n"+
					"%s\n\n"+
					"The expiration date is 30 minutes after issuance. \n\n"+
					"If you don't know, please discard this email. Unless you click on the URL, it will not be processed.",
				fmt.Sprintf("%s/account/forgot/%s", config.Config.Server.Endpoint, verifyToken),
			)
			err := SendMail(targetEmail, "Forgot password", "text/plain", body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateForgotPasswordToken(token string) (uuid.UUID, error) {
	var targetUuidRaw string
	err := database.DBPool.QueryRow(
		context.Background(),
		"select target from user_reset_password where key = $1 and expired > now();",
		token,
	).Scan(
		&targetUuidRaw,
	)
	if err != nil {
		return uuid.Nil, err
	}

	targetUuid := uuid.MustParse(targetUuidRaw)
	return targetUuid, nil
}

func IssueVerifyEmail(targetUuid uuid.UUID, targetEmail string, tx ...pgx.Tx) error {
	var issueSql = "insert into user_mail_transaction(uuid, new_email, token) VALUES ($1, $2, $3) on conflict on constraint user_mail_transaction_pk do update set new_email = $2, token = $3, expired_at = default;"
	verifyToken := misc.GenToken(32)

	var err error

	if len(tx) > 0 {
		_, err = tx[0].Exec(
			context.Background(),
			issueSql,
			targetUuid.String(),
			targetEmail,
			verifyToken,
		)
	} else {
		_, err = database.DBPool.Exec(
			context.Background(),
			issueSql,
			targetUuid.String(),
			targetEmail,
			verifyToken,
		)
	}
	if err != nil {
		return err
	}

	// TODO: 将来的に"template/text"を使う
	body := fmt.Sprintf(
		"Hello. This is MatticNote.\n\n"+
			"I was asked to use this email address to verify my account email and issued a token to accept it.\n\n"+
			"Click the URL below to complete the email verification: \n\n"+
			"%s\n\n"+
			"The expiration date is 3 hours after issuance. \n\n"+
			"If you don't know, please discard this email. Unless you click on the URL, it will not be accepted.",
		fmt.Sprintf("%s/account/verify/%s", config.Config.Server.Endpoint, verifyToken),
	)

	if err := SendMail(targetEmail, "Verify Email", "text/plain", body); err != nil {
		return err
	}

	return nil
}

func ChangeUserPassword(targetUuid uuid.UUID, password string) error {
	var email string

	err := database.DBPool.QueryRow(
		context.Background(),
		"select email from user_mail where uuid = $1 and is_verified is true;",
		targetUuid.String(),
	).Scan(
		&email,
	)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), UserPasswordHashCost)
	if err != nil {
		return err
	}

	dbRes, err := database.DBPool.Exec(
		context.Background(),
		"update user_password set password = $1 where uuid = $2;\n",
		hashedPassword,
		targetUuid,
	)
	if err != nil {
		return err
	}

	if dbRes.RowsAffected() > 0 {
		// TODO: 将来的に"template/text"を使う
		body := "Hello. This is MatticNote.\n\n" +
			"We would like to inform you that the password of the account associated with this email address has been changed recently.\n\n" +
			"If you did it yourself, there is no problem. If you don't remember making any changes, please let the instance administrator know immediately."

		if err := SendMail(email, "Password was changed", "text/plain", body); err != nil {
			return err
		}
	}

	return nil
}

func UpdateProfile(targetUuid uuid.UUID, name, summary string, isBot, acceptManually bool) error {
	_, err := database.DBPool.Exec(
		context.Background(),
		"update \"user\" set display_name = $1, summary = $2, is_bot = $3, accept_manually = $4, updated_at = now() where uuid = $5;",
		name,
		summary,
		isBot,
		acceptManually,
		targetUuid.String(),
	)

	return err
}

func InsertSigninLog(targetUuid uuid.UUID, fromIp string, isSuccess bool) (err error) {
	if targetUuid != uuid.Nil {
		_, err = database.DBPool.Exec(
			context.Background(),
			"insert into signin(target_user, from_ip, is_success) values ($1, $2, $3);",
			targetUuid.String(),
			fromIp,
			isSuccess,
		)
	} else {
		_, err = database.DBPool.Exec(
			context.Background(),
			"insert into signin(from_ip, is_success) values ($1, $2);",
			fromIp,
			isSuccess,
		)
	}

	return err
}
