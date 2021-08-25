package user

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/ist"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"strings"
)

//goland:noinspection GoUnusedGlobalVariable

const (
	PasswordHashCost = 12
	KeyPairLength    = 2048
)

func RegisterLocalUser(email, username, password string, skipEmailVerify bool) (*uuid.UUID, error) {
	var count int
	err := database.DBPool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM \"user\" LEFT JOIN user_mail um on \"user\".uuid = um.uuid WHERE (username ILIKE $1 AND host IS NULL) OR email ILIKE $2;",
		username,
		email,
	).Scan(&count)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if count > 0 {
		return nil, ErrUserExists
	}

	newUuid := uuid.Must(uuid.NewRandom())

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return nil, err
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
		return nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_mail(uuid, email, is_verified) VALUES ($1, $2, $3);",
		newUuid,
		email,
		skipEmailVerify,
	)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), PasswordHashCost)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_password(uuid, password) VALUES ($1, $2);",
		newUuid,
		hashedPassword,
	)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, rsaPublicKey := misc.GenerateRSAKeypair(KeyPairLength)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_signature_key(uuid, public_key, private_key) VALUES ($1, $2, $3);",
		newUuid,
		string(rsaPublicKey),
		string(rsaPrivateKey),
	)
	if err != nil {
		return nil, err
	}

	if !skipEmailVerify {
		if err := IssueVerifyEmail(newUuid, email, tx); err != nil {
			return nil, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &newUuid, err
}

func GetUser(targetUuid uuid.UUID) (*ist.UserStruct, error) {
	target := new(ist.UserStruct)
	var (
		isActive     bool
		publicKeyRaw string
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select \"user\".uuid, host, username, display_name, summary, created_at, updated_at, is_silence, accept_manually, is_bot, is_suspend, is_active, public_key "+
			"from \"user\" left join user_signature_key usk on \"user\".uuid = usk.uuid where \"user\".uuid = $1;",
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
		&target.IsSuspend,
		&isActive,
		&publicKeyRaw,
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

	target.PublicKey, _ = pem.Decode([]byte(publicKeyRaw))

	return target, nil
}

func GetLocalUser(targetUuid uuid.UUID) (*ist.LocalUserStruct, error) {
	target := new(ist.LocalUserStruct)
	var (
		isActive      bool
		publicKeyRaw  string
		privateKeyRaw string
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select \"user\".uuid, username, email, display_name, summary, created_at, updated_at, is_silence, accept_manually, is_superuser, is_bot, is_suspend, is_active, public_key, private_key "+
			"from \"user\" left join user_mail um on \"user\".uuid = um.uuid left join user_signature_key usk on \"user\".uuid = usk.uuid where \"user\".uuid = $1 and \"user\".host is null",
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
		&target.IsSuspend,
		&isActive,
		&publicKeyRaw,
		&privateKeyRaw,
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

	target.PublicKey, _ = pem.Decode([]byte(publicKeyRaw))
	if target.PublicKey.Type != misc.PublicKeyType {
		panic(fmt.Sprintf("public key type is not \"%s\"", misc.PublicKeyType))
	}
	target.PrivateKey, _ = pem.Decode([]byte(privateKeyRaw))
	if target.PrivateKey.Type != misc.PrivateKeyType {
		panic(fmt.Sprintf("private key type is not \"%s\"", misc.PublicKeyType))
	}

	return target, nil
}

func GetLocalUserFromUsername(username string) (*ist.LocalUserStruct, error) {
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

func LookupUserRelation(fromUuid, targetUuid uuid.UUID) (*ist.UserRelationStruct, error) {
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

	relationStruct := new(ist.UserRelationStruct)

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

func GetUserPublicKeyFromKeyId(keyId string) (*pem.Block, error) {
	var (
		err          error
		publicKeyRaw string
		isActive     bool
		isSuspend    bool
	)

	if strings.HasPrefix(keyId, config.Config.Server.Endpoint) {
		urlParse, err := url.Parse(keyId)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(urlParse.Path, "/activity/user/") {
			return nil, ErrInvalidKeyId
		}
		targetUuidStr := strings.Replace(
			strings.TrimPrefix(urlParse.Path, "/activity/user/"), "/", "", -1,
		)
		targetUuid, err := uuid.Parse(targetUuidStr)
		if err != nil {
			return nil, ErrInvalidKeyId
		}
		err = database.DBPool.QueryRow(
			context.Background(),
			"select public_key, is_active, is_suspend from user_signature_key, \"user\" where \"user\".uuid = $1 and \"user\".uuid = user_signature_key.uuid;",
			targetUuid.String(),
		).Scan(
			&publicKeyRaw,
			&isActive,
			&isSuspend,
		)
	} else {
		err = database.DBPool.QueryRow(
			context.Background(),
			"select public_key, is_active, is_suspend from user_signature_key, \"user\" where key_id = $1 and \"user\".uuid = user_signature_key.uuid;",
			keyId,
		).Scan(
			&publicKeyRaw,
			&isActive,
			&isSuspend,
		)
	}

	if err != nil && err == pgx.ErrNoRows {
		return nil, ErrNoSuchUser
	} else if err != nil {
		return nil, err
	}

	if !isActive {
		return nil, ErrUserGone
	}

	if publicKeyRaw == "" {
		return nil, ErrInvalidKey
	}

	pubKeyPem, _ := pem.Decode([]byte(publicKeyRaw))

	if pubKeyPem == nil || pubKeyPem.Type != misc.PublicKeyType {
		return nil, ErrInvalidKey
	}

	return pubKeyPem, nil
}

func GetUserPrivateKey(targetUuid uuid.UUID) (*pem.Block, error) {
	var (
		privateKeyRaw string
		isActive      bool
		isSuspend     bool
	)
	err := database.DBPool.QueryRow(
		context.Background(),
		"select private_key, is_active, is_suspend from user_signature_key, \"user\" where \"user\".uuid = $1 and \"user\".uuid = user_signature_key.uuid and \"user\".host is null;",
		targetUuid.String(),
	).Scan(
		&privateKeyRaw,
		&isActive,
		&isSuspend,
	)
	if err != nil && err == pgx.ErrNoRows {
		return nil, ErrNoSuchUser
	} else if err != nil {
		return nil, err
	}

	if !isActive {
		return nil, ErrUserGone
	}

	if isSuspend {
		return nil, ErrUserSuspended
	}

	if privateKeyRaw == "" {
		return nil, ErrInvalidKey
	}

	privateKey, _ := pem.Decode([]byte(privateKeyRaw))

	if privateKey == nil || privateKey.Type != misc.PrivateKeyType {
		return nil, ErrInvalidKey
	}

	return privateKey, nil
}
