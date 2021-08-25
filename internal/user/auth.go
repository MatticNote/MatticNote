package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/mail"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/mn_template"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func ValidateLoginUser(login, password string) (uuid.UUID, error) {
	var (
		targetUuid     uuid.UUID
		isMailVerified bool
		targetPassword []byte
		is2faEnabled   misc.NullBool
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"SELECT \"user\".uuid, is_verified, password, is_enable FROM \"user\" LEFT JOIN user_2fa u2f ON \"user\".uuid = u2f.uuid JOIN user_mail um ON \"user\".uuid = um.uuid JOIN user_password up on \"user\".uuid = up.uuid WHERE email ILIKE $1 OR (username ILIKE $1 AND host IS NULL);",
		login,
	).Scan(
		&targetUuid,
		&isMailVerified,
		&targetPassword,
		&is2faEnabled,
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

	if is2faEnabled.Valid && is2faEnabled.Bool {
		return targetUuid, Err2faRequired
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
			txtTemplate, err := mn_template.LoadTextTemplate()
			if err != nil {
				return err
			}

			body := new(bytes.Buffer)

			if err := txtTemplate.ExecuteTemplate(body, "issue_password.txt", fiber.Map{
				"ResetURL": fmt.Sprintf("%s/account/forgot/%s", config.Config.Server.Endpoint, verifyToken),
			}); err != nil {
				return err
			}

			if err := mail.SendMail(targetEmail, "Forgot password", "text/plain", body.String()); err != nil {
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

	body := new(bytes.Buffer)

	txtTemplate, err := mn_template.LoadTextTemplate()
	if err != nil {
		return err
	}

	if err := txtTemplate.ExecuteTemplate(body, "verify_mail.txt", fiber.Map{
		"VerifyURL": fmt.Sprintf("%s/account/verify/%s", config.Config.Server.Endpoint, verifyToken),
	}); err != nil {
		return err
	}

	if err := mail.SendMail(targetEmail, "Verify Email", "text/plain", body.String()); err != nil {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), PasswordHashCost)
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
		body := new(bytes.Buffer)

		txtTemplate, err := mn_template.LoadTextTemplate()
		if err != nil {
			return err
		}

		if err := txtTemplate.ExecuteTemplate(body, "password_changed.txt", nil); err != nil {
			return err
		}

		if err := mail.SendMail(email, "Password was changed", "text/plain", body.String()); err != nil {
			return err
		}
	}

	return nil
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

func Setup2faCode(targetUUid uuid.UUID) (*otp.Key, error) {
	user, err := GetLocalUser(targetUUid)
	if err != nil {
		return nil, err
	}

	var (
		isEnable   bool
		totpGen    *otp.Key
		secretCode string
	)

	err = database.DBPool.QueryRow(
		context.Background(),
		"select is_enable, secret_code from user_2fa where uuid = $1;",
		user.Uuid.String(),
	).Scan(
		&isEnable,
		&secretCode,
	)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}

		secret := make([]byte, 20)
		_, err := rand.Reader.Read(secret)
		if err != nil {
			return nil, err
		}
		secretCode = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

		backupCode, err := json.Marshal(misc.GenBackupCode())
		if err != nil {
			return nil, err
		}

		_, err = database.DBPool.Exec(
			context.Background(),
			"insert into user_2fa(uuid, secret_code, backup_code) values ($1, $2, $3);",
			user.Uuid.String(),
			secretCode,
			backupCode,
		)
		if err != nil {
			return nil, err
		}
	}

	if isEnable {
		return nil, Err2faAlreadyEnabled
	}

	totpGen, err = otp.NewKeyFromURL(fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s&period=30&digits=6&algorithm=SHA1",
		fmt.Sprintf("MatticNote(%s)", config.Config.Server.Endpoint),
		user.Username,
		secretCode,
		fmt.Sprintf("MatticNote(%s)", config.Config.Server.Endpoint),
	))
	if err != nil {
		return nil, err
	}

	return totpGen, nil
}

func Validate2faCode(targetUuid uuid.UUID, token string) error {
	var secretCode string
	err := database.DBPool.QueryRow(
		context.Background(),
		"select secret_code from user_2fa where uuid = $1;",
		targetUuid.String(),
	).Scan(
		&secretCode,
	)
	if err != nil {
		return err
	}

	if !totp.Validate(token, secretCode) {
		return ErrInvalid2faToken
	}

	return nil
}

func Enable2faAuth(targetUuid uuid.UUID) error {
	exec, err := database.DBPool.Exec(
		context.Background(),
		"update user_2fa set is_enable = true where uuid = $1;",
		targetUuid.String(),
	)
	if err != nil {
		return err
	}

	if exec.RowsAffected() <= 0 {
		return ErrCantEnable2fa
	}

	return nil
}

func StatusUser2fa(targetUuid uuid.UUID) (bool, error) {
	var isEnabled bool
	err := database.DBPool.QueryRow(
		context.Background(),
		"select count(*) > 0 as is_enabled from user_2fa where uuid = $1;",
		targetUuid.String(),
	).Scan(
		&isEnabled,
	)

	return isEnabled, err
}

func Disable2faAuth(targetUuid uuid.UUID) error {
	exec, err := database.DBPool.Exec(
		context.Background(),
		"delete from user_2fa where uuid = $1;",
		targetUuid.String(),
	)
	if err != nil {
		return err
	}

	if exec.RowsAffected() <= 0 {
		return ErrCantDisable2fa
	}

	return nil
}

func Get2faBackupCode(targetUuid uuid.UUID) (code []string, err error) {
	err = database.DBPool.QueryRow(
		context.Background(),
		"select backup_code from user_2fa where uuid = $1 and is_enable = true;",
		targetUuid.String(),
	).Scan(
		&code,
	)

	return code, err
}

func Use2faBackupCode(targetUuid uuid.UUID, code string) error {
	currentCode, err := Get2faBackupCode(targetUuid)
	if err != nil {
		return err
	}

	for i, v := range currentCode {
		if code == v {
			currentCode[i] = ""
			currentCodeJson, err := json.Marshal(currentCode)
			if err != nil {
				return err
			}
			_, err = database.DBPool.Exec(
				context.Background(),
				"update user_2fa set backup_code = $1 where uuid = $2;",
				currentCodeJson,
				targetUuid.String(),
			)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return ErrInvalid2faToken
}

func Regenerate2faBackupCode(targetUuid uuid.UUID) error {
	newBackupCode, err := json.Marshal(misc.GenBackupCode())
	if err != nil {
		return err
	}

	_, err = database.DBPool.Exec(
		context.Background(),
		"update user_2fa set backup_code = $1 where uuid = $2;",
		newBackupCode,
		targetUuid.String(),
	)
	if err != nil {
		return err
	}

	return nil
}
