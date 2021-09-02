package oauth

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/ory/fosite"
	"time"
)

//var store *storage.MemoryStore
var store = &MNOAuthStore{}

type MNOAuthStore struct {
}

func (s *MNOAuthStore) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	client := new(MNOAuthClient)
	err := database.DBPool.QueryRow(
		ctx,
		"select client_key, client_secret, name, client_owner, redirect_uris, scopes from oauth_client where client_key = $1;",
		id,
	).Scan(
		&client.ID,
		&client.Secret,
		&client.Name,
		&client.Owner,
		&client.RedirectURIs,
		&client.Scopes,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fosite.ErrNotFound
		} else {
			return nil, err
		}
	}

	return client, nil
}

func (s *MNOAuthStore) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	var expiredAt time.Time
	err := database.DBPool.QueryRow(
		ctx,
		"select expired_at from oauth_valid_jti where jti = $1;",
		jti,
	).Scan(
		&expiredAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	if expiredAt.After(time.Now()) {
		return fosite.ErrJTIKnown
	} else {
		return nil
	}
}

func (s *MNOAuthStore) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	_, err := database.DBPool.Exec(
		ctx,
		"delete from oauth_valid_jti where expired_at < now();",
	)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}
	exec, err := database.DBPool.Exec(
		ctx,
		"insert into oauth_valid_jti(jti, expired_at) values ($1, $2) on conflict do nothing;",
		jti,
		exp,
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return fosite.ErrJTIKnown
	}

	return nil
}

func (s *MNOAuthStore) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {
	parsedUserUuid, err := uuid.Parse(req.GetSession().GetUsername())
	if err != nil {
		return errors.New("username is not uuid")
	}
	_, err = database.DBPool.Exec(
		ctx,
		"insert into oauth_authorize_code(code, expires_at, scopes, client_id, user_id) values ($1, $2, $3, $4, $5);",
		code,
		req.GetSession().GetExpiresAt(fosite.AuthorizeCode),
		req.GetGrantedScopes(),
		req.GetClient().GetID(),
		parsedUserUuid.String(),
	)
	return err
}

func (s *MNOAuthStore) GetAuthorizeCodeSession(ctx context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	var (
		userId    string
		expiresAt time.Time
		clientId  string
		scopes    fosite.Arguments
		isActive  bool
	)
	err := database.DBPool.QueryRow(
		ctx,
		"select expires_at, client_id, scopes, user_id, is_active from oauth_authorize_code where code = $1;",
		code,
	).Scan(
		&expiresAt,
		&clientId,
		&scopes,
		&userId,
		&isActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fosite.ErrNotFound
		} else {
			return nil, err
		}
	}
	req := fosite.NewRequest()
	req.GrantedScope = scopes
	req.SetSession(&fosite.DefaultSession{
		Username: userId,
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AuthorizeCode: expiresAt,
		},
	})

	cli, err := s.GetClient(ctx, clientId)
	if err != nil {
		return nil, err
	}
	req.Client = cli

	if !isActive {
		return req, fosite.ErrInvalidatedAuthorizeCode
	}

	return req, nil
}

func (s *MNOAuthStore) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	_, err := database.DBPool.Exec(
		ctx,
		"update oauth_authorize_code set is_active = false where code = $1;",
		code,
	)
	return err
}

func (s *MNOAuthStore) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	userUuid, err := uuid.Parse(req.GetSession().GetUsername())
	if err != nil {
		return err
	}

	_, err = database.DBPool.Exec(
		ctx,
		"insert into oauth_access_token(token, req_id, expires, scopes, client_id, user_id) values ($1, $2, $3, $4, $5, $6);",
		signature,
		req.GetID(),
		req.GetSession().GetExpiresAt(fosite.AccessToken),
		req.GetGrantedScopes(),
		req.GetClient().GetID(),
		userUuid.String(),
	)
	return err
}

func (s *MNOAuthStore) GetAccessTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	var (
		reqId    string
		expires  time.Time
		scopes   fosite.Arguments
		clientId string
		userId   uuid.UUID
	)

	err := database.DBPool.QueryRow(
		ctx,
		"select req_id, expires, scopes, client_id, user_id from oauth_access_token where token = $1;",
		signature,
	).Scan(
		&reqId,
		&expires,
		&scopes,
		&clientId,
		&userId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fosite.ErrNotFound
		} else {
			return nil, err
		}
	}
	req := fosite.NewRequest()
	req.SetSession(&fosite.DefaultSession{
		Username: userId.String(),
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken: expires,
		},
	})
	req.SetID(reqId)
	req.GrantedScope = scopes
	cli, err := s.GetClient(ctx, clientId)
	if err != nil {
		return nil, err
	}
	req.Client = cli

	return req, nil
}

func (s *MNOAuthStore) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	_, err := database.DBPool.Exec(
		ctx,
		"delete from oauth_access_token where token = $1;",
		signature,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MNOAuthStore) CreateRefreshTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	userUuid, err := uuid.Parse(req.GetSession().GetUsername())
	if err != nil {
		return err
	}

	_, err = database.DBPool.Exec(
		ctx,
		"insert into oauth_refresh_token(token, req_id, client_id, user_id, expires_at, scopes) values ($1, $2, $3, $4, $5, $6);",
		signature,
		req.GetID(),
		req.GetClient().GetID(),
		userUuid.String(),
		func() interface{} {
			expiresAt := req.GetSession().GetExpiresAt(fosite.RefreshToken)
			if expiresAt.IsZero() {
				return nil
			} else {
				return expiresAt
			}
		}(),
		req.GetGrantedScopes(),
	)
	return err
}

func (s *MNOAuthStore) GetRefreshTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	var (
		reqId    string
		expires  misc.NullTime
		scopes   fosite.Arguments
		clientId string
		userId   uuid.UUID
	)
	err := database.DBPool.QueryRow(
		ctx,
		"select req_id, user_id, client_id, scopes, expires_at from oauth_refresh_token where token = $1;",
		signature,
	).Scan(
		&reqId,
		&userId,
		&clientId,
		&scopes,
		&expires,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fosite.ErrNotFound
		} else {
			return nil, err
		}
	}

	req := fosite.NewRequest()
	req.SetSession(&fosite.DefaultSession{
		Username: userId.String(),
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.RefreshToken: func() time.Time {
				if expires.Valid {
					return expires.Time
				} else {
					return time.Time{}
				}
			}(),
		},
	})
	req.SetID(reqId)
	req.GrantedScope = scopes
	cli, err := s.GetClient(ctx, clientId)
	if err != nil {
		return nil, err
	}
	req.Client = cli

	return req, nil
}

func (s *MNOAuthStore) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	_, err := database.DBPool.Exec(
		ctx,
		"delete from oauth_refresh_token where token = $1;",
		signature,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MNOAuthStore) RevokeRefreshToken(ctx context.Context, requestID string) error {
	_, err := database.DBPool.Exec(
		ctx,
		"delete from oauth_refresh_token where req_id = $1;",
		requestID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MNOAuthStore) RevokeAccessToken(ctx context.Context, requestID string) error {
	_, err := database.DBPool.Exec(
		ctx,
		"delete from oauth_access_token where req_id = $1;",
		requestID,
	)
	if err != nil {
		return err
	}

	return nil
}

//func (s *MNOAuthStore) GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (*jose.JSONWebKey, error) {
//	return nil, nil
//}
//
//func (s *MNOAuthStore) GetPublicKeys(ctx context.Context, issuer string, subject string) (*jose.JSONWebKeySet, error) {
//	return nil, nil
//}
//
//func (s *MNOAuthStore) GetPublicKeyScopes(ctx context.Context, issuer string, subject string, keyId string) ([]string, error) {
//	return nil, nil
//}
//
//func (s *MNOAuthStore) IsJWTUsed(ctx context.Context, jti string) (bool, error) {
//	return false, nil
//}
//
//func (s *MNOAuthStore) MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) error {
//	return nil
//}
