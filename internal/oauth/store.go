package oauth

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/jackc/pgx/v4"
	"github.com/ory/fosite"
	"github.com/ory/fosite/storage"
	"time"
)

var store *storage.MemoryStore

type mnOAuthStore struct {
}

func (s *mnOAuthStore) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	client := new(mnOAuthClient)
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

func (s *mnOAuthStore) ClientAssertionJWTValid(ctx context.Context, jti string) error {
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

func (s *mnOAuthStore) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
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
