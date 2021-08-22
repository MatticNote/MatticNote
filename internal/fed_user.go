package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/piprate/json-gold/ld"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrInvalidResponse     = errors.New("invalid response")
	ErrUnknownASType       = errors.New("unknown activitystreams type")
	ErrUnknownASId         = errors.New("unknown activitystreams id")
	ErrInvalidValueType    = errors.New("invalid value type")
	ErrInvalidIdUrl        = errors.New("invalid id url")
	ErrNotEnoughParams     = errors.New("not enough params")
	ErrKeyOwnerDoesntMatch = errors.New("key owner does not match")
	ErrFailedFetch         = errors.New("failed fetch remote user")
)

var (
	jldProc    = ld.NewJsonLdProcessor()
	jldOptions = ld.NewJsonLdOptions("")
	jldDoc     = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		map[string]interface{}{
			"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
			"toot":                      "http://joinmastodon.org/ns#",
			"featured": map[string]interface{}{
				"@id":   "toot:featured",
				"@type": "@id",
			},
			"alsoKnownAs": map[string]interface{}{
				"@id":   "toot:alsoKnownAs",
				"@type": "@id",
			},
			"movedTo": map[string]interface{}{
				"@id":   "toot:movedTo",
				"@type": "@id",
			},
			"schema":        "http://schema.org#",
			"PropertyValue": "schema:PropertyValue",
			"value":         "schema:value",
			"discoverable":  "toot:discoverable",
		},
	}
)

func RegisterRemoteUser(actorUrl string) (*uuid.UUID, error) {
	req, err := http.NewRequest(http.MethodGet, actorUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/activity+json")
	req.Header.Set("User-Agent", fmt.Sprintf("MatticNote/%s", Version))

	// todo: 将来的にはHTTP Signatureとかを付ける処理をここでやる

	cli := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode >= 300 {
		return nil, ErrFailedFetch
	}

	bufBody := new(bytes.Buffer)
	_, err = bufBody.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}

	var apData map[string]interface{}
	if err := json.Unmarshal(bufBody.Bytes(), &apData); err != nil {
		return nil, ErrInvalidResponse
	}

	apDataDoc, err := jldProc.Compact(apData, jldDoc, jldOptions)
	if err != nil {
		return nil, err
	}

	apType, ok := apDataDoc["type"]
	if !ok {
		return nil, ErrUnknownASType
	}
	var isBot = false
	switch apType {
	case "Person":
		// Valid
	case "Service":
		// Valid
		isBot = true
	default:
		// Invalid
		return nil, ErrUnknownASType
	}

	apId, ok := apDataDoc["id"]
	if !ok {
		return nil, ErrUnknownASId
	}
	apIdStr, ok := apId.(string)
	if !ok {
		return nil, ErrInvalidValueType
	}
	apIdParsed, err := url.Parse(apIdStr)
	if err != nil {
		return nil, ErrInvalidIdUrl
	}

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	var sharedInboxP *string
	endpointsRaw, ok := apDataDoc["endpoints"]
	if ok {
		endpoints, ok := endpointsRaw.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidValueType
		}
		sharedInboxRaw, ok := endpoints["sharedInbox"]
		if ok {
			sharedInbox, ok := sharedInboxRaw.(string)
			if !ok {
				return nil, ErrInvalidValueType
			}
			sharedInboxP = &sharedInbox
		}
	}

	_, err = tx.Exec(
		context.Background(),
		"insert into host(host, shared_inbox) values ($1, $2) on conflict do nothing;",
		apIdParsed.Host,
		sharedInboxP,
	)
	if err != nil {
		return nil, err
	}

	newUuid := uuid.Must(uuid.NewRandom())

	usernameRaw, ok := apDataDoc["preferredUsername"]
	if !ok {
		return nil, ErrNotEnoughParams
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, ErrInvalidValueType
	}
	var displayNameP *string = nil
	displayNameRaw, ok := apDataDoc["name"]
	if ok {
		displayName, ok := displayNameRaw.(string)
		if !ok {
			return nil, ErrInvalidValueType
		} else {
			displayNameP = &displayName
		}
	}
	var summaryP *string = nil
	summaryRaw, ok := apDataDoc["summary"]
	if ok {
		summary, ok := summaryRaw.(string)
		if !ok {
			return nil, ErrInvalidValueType
		} else {
			summaryP = &summary
		}
	}
	var acceptManually = false
	acceptManuallyRaw, ok := apDataDoc["manuallyApprovesFollowers"]
	if ok {
		acceptManually, ok = acceptManuallyRaw.(bool)
		if !ok {
			return nil, ErrInvalidValueType
		}
	}

	_, err = tx.Exec(
		context.Background(),
		"insert into \"user\"(uuid, username, host, display_name, summary, accept_manually, is_bot, ap_id) values ($1, $2, $3, $4, $5, $6, $7, $8);",
		newUuid,
		username,
		apIdParsed.Host,
		displayNameP,
		summaryP,
		acceptManually,
		isBot,
		apIdStr,
	)
	if err != nil {
		return nil, err
	}

	publicKeyRaw, ok := apDataDoc["publicKey"]
	if ok {
		publicKeyMI, ok := publicKeyRaw.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidValueType
		}
		keyIdRaw, ok := publicKeyMI["id"]
		if !ok {
			return nil, ErrNotEnoughParams
		}
		keyId, ok := keyIdRaw.(string)
		if !ok {
			return nil, ErrInvalidValueType
		}
		keyOwnerRaw, ok := publicKeyMI["owner"]
		if !ok {
			return nil, ErrNotEnoughParams
		}
		keyOwner, ok := keyOwnerRaw.(string)
		if !ok {
			return nil, ErrInvalidValueType
		}
		if keyOwner != apIdStr {
			return nil, ErrKeyOwnerDoesntMatch
		}
		publicKeyPemRaw, ok := publicKeyMI["publicKeyPem"]
		if !ok {
			return nil, ErrNotEnoughParams
		}
		publicKeyPem, ok := publicKeyPemRaw.(string)
		if !ok {
			return nil, ErrInvalidValueType
		}
		_, err = tx.Exec(
			context.Background(),
			"insert into user_signature_key(uuid, public_key, key_id) values ($1, $2, $3);",
			newUuid,
			publicKeyPem,
			keyId,
		)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &newUuid, nil
}

func GetRemoteUserFromApID(apId string) (*UserStruct, error) {
	var targetUuid uuid.UUID
	err := database.DBPool.QueryRow(
		context.Background(),
		"select uuid from \"user\" where ap_id = $1",
		apId,
	).Scan(
		&targetUuid,
	)
	if err != nil && err == pgx.ErrNoRows {
		return nil, ErrNoSuchUser
	} else if err != nil {
		return nil, err
	}

	return GetUser(targetUuid)
}
