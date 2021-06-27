package internal

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/mn_template"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrUserExists        = errors.New("username or email is already taken")
	ErrLoginFailed       = errors.New("there is an error in the login name or password")
	ErrEmailAuthRequired = errors.New("target user required email authentication")
	Err2faRequired       = errors.New("target user required two factor authentication")
	ErrNoSuchUser        = errors.New("target user was not found")
	ErrUserGone          = errors.New("target user was gone")
	ErrUserSuspended     = errors.New("target user is suspended")
	ErrInvalidPassword   = errors.New("invalid password")
	Err2faAlreadyEnabled = errors.New("target user 2fa is already enabled")
	ErrInvalid2faToken   = errors.New("invalid 2fa token")
	ErrCantEnable2fa     = errors.New("cannot enable target user's 2fa")
)

type (
	LocalUserStruct struct {
		UserStruct
		Email          string
		AcceptManually bool
		IsSuperuser    bool
	}
	UserStruct struct {
		Uuid           uuid.UUID
		Username       string
		Host           misc.NullString
		DisplayName    misc.NullString
		Summary        misc.NullString
		CreatedAt      misc.NullTime
		UpdatedAt      misc.NullTime
		IsSilence      bool
		AcceptManually bool
		IsBot          bool
	}
	UserRelationStruct struct {
		Following     bool
		FollowPending bool
		Follows       bool
		Muting        bool
		Blocking      bool
		Blocked       bool
	}
)

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
		if err := IssueVerifyEmail(newUuid, email, tx); err != nil {
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

func GetUser(targetUuid uuid.UUID) (*UserStruct, error) {
	target := new(UserStruct)
	var (
		isSuspend bool
		isActive  bool
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select uuid, host, username, display_name, summary, created_at, updated_at, is_silence, accept_manually, is_bot, is_suspend, is_active from \"user\" where \"user\".uuid = $1",
		targetUuid.String(),
	).Scan(
		&target.Uuid,
		&target.Host,
		&target.Username,
		&target.DisplayName,
		&target.Summary,
		&target.CreatedAt,
		&target.UpdatedAt,
		&target.IsSilence,
		&target.AcceptManually,
		&target.IsBot,
		&isSuspend,
		&isActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoSuchUser
		} else {
			return nil, err
		}
	}

	if !isActive {
		return nil, ErrUserGone
	}

	if isSuspend {
		return nil, ErrUserSuspended
	}

	return target, nil
}

func GetLocalUser(targetUuid uuid.UUID) (*LocalUserStruct, error) {
	target := new(LocalUserStruct)
	var (
		isSuspend bool
		isActive  bool
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select \"user\".uuid, username, email, display_name, summary, created_at, updated_at, is_silence, accept_manually, is_superuser, is_bot, is_suspend, is_active from \"user\" left join user_mail um on \"user\".uuid = um.uuid where \"user\".uuid = $1",
		targetUuid.String(),
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
		&isActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoSuchUser
		} else {
			return nil, err
		}
	}

	if !isActive {
		return nil, ErrUserGone
	}

	if isSuspend {
		return nil, ErrUserSuspended
	}

	return target, nil
}

func GetLocalUserFromUsername(username string) (*LocalUserStruct, error) {
	var targetUuid uuid.UUID

	err := database.DBPool.QueryRow(
		context.Background(),
		"select uuid from \"user\" where username ilike $1 and host is null",
		username,
	).Scan(
		&targetUuid,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoSuchUser
		} else {
			return nil, err
		}
	}

	return GetLocalUser(targetUuid)
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

			if err := SendMail(targetEmail, "Forgot password", "text/plain", body.String()); err != nil {
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

	if err := SendMail(targetEmail, "Verify Email", "text/plain", body.String()); err != nil {
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
		body := new(bytes.Buffer)

		txtTemplate, err := mn_template.LoadTextTemplate()
		if err != nil {
			return err
		}

		if err := txtTemplate.ExecuteTemplate(body, "password_changed.txt", nil); err != nil {
			return err
		}

		if err := SendMail(email, "Password was changed", "text/plain", body.String()); err != nil {
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

func CheckUsernameUsed(username string) (bool, error) {
	var checkCount int
	err := database.DBPool.QueryRow(
		context.Background(),
		"select count(*) from \"user\" where username ilike $1 and host is null;",
		username,
	).Scan(
		&checkCount,
	)
	if err != nil {
		return false, err
	}

	return checkCount > 0, nil
}

func LookupUserRelation(fromUuid, targetUuid uuid.UUID) (*UserRelationStruct, error) {
	if fromUuid == targetUuid {
		return nil, ErrCantRelateYourself
	}

	rows, err := database.DBPool.Query(
		context.Background(),
		"select 'following' as relation from follow_relation where follow_from = $1 and follow_to = $2 and is_pending = false "+
			"union select 'follow_pending' from follow_relation where follow_from = $1 and follow_to = $2 and is_pending = true "+
			"union select 'follows' from follow_relation where follow_from = $1 and follow_to = $2 and is_pending = false "+
			"union select 'muting' from mute_relation where mute_from = $1 and mute_to = $2"+
			"union select 'blocking' from block_relation where block_from = $1 and block_to = $2 "+
			"union select 'blocked' from block_relation where block_to = $1 and block_from = $2;",
		fromUuid.String(),
		targetUuid.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relationStruct := new(UserRelationStruct)

	for rows.Next() {
		var relation string
		err := rows.Scan(&relation)
		if err != nil {
			return nil, err
		}
		switch relation {
		case "following":
			relationStruct.Following = true
		case "follow_pending":
			relationStruct.FollowPending = true
		case "follows":
			relationStruct.Follows = true
		case "muting":
			relationStruct.Muting = true
		case "blocking":
			relationStruct.Blocking = true
		case "blocked":
			relationStruct.Blocked = true
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return relationStruct, nil
}
