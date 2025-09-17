package repositories

import (
	"context"
	"time"

	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	"github.com/jacoobjake/einvoice-api/internal/database/enums"
	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/pkg/errors"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

var AuthTokens = models.AuthTokens

type AuthTokenRepository struct {
	db bob.Executor
}

func (r *AuthTokenRepository) Create(ctx context.Context, token *models.AuthTokenSetter) (*models.AuthToken, error) {
	createdToken, err := AuthTokens.Insert(token).One(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting auth_tokens")
	}
	return createdToken, nil
}

func (r *AuthTokenRepository) FindByToken(ctx context.Context, token string) (*models.AuthToken, error) {
	authToken, err := AuthTokens.Query(
		sm.Where(AuthTokens.Columns.Token.EQ(psql.Arg(token))),
	).One(ctx, r.db)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching token")
	}

	return authToken, nil
}

func (r *AuthTokenRepository) FindTokenByUserIdAndType(ctx context.Context, userID int64, tokenType enums.AuthTokenTypes) (*models.AuthToken, error) {
	authToken, err := AuthTokens.Query(
		sm.Where(AuthTokens.Columns.UserID.EQ(psql.Arg(userID))),
		sm.Where(AuthTokens.Columns.Type.EQ(psql.Arg(tokenType))),
	).One(ctx, r.db)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching token")
	}

	return authToken, nil
}

func (r *AuthTokenRepository) InvalidateActiveTokensByUserID(ctx context.Context, userID int64, tokenType enums.AuthTokenTypes) error {
	invalidate := models.AuthTokenSetter{
		ExpireAt: omitnull.From(time.Now()),
	}

	_, err := AuthTokens.Update(
		invalidate.UpdateMod(),
		um.Where(
			psql.And(
				AuthTokens.Columns.UserID.EQ(psql.Arg(userID)),
				AuthTokens.Columns.Type.EQ(psql.Arg(tokenType)),
				AuthTokens.Columns.ExpireAt.GTE(psql.Arg(time.Now())),
			),
		),
	).All(ctx, r.db)

	if err != nil {
		return errors.Wrap(err, "error executing invalidate token query")
	}

	return nil
}

func (r *AuthTokenRepository) InvalidateActiveTokensBySessionID(ctx context.Context, sessionID uuid.UUID, tokenType enums.AuthTokenTypes) error {
	invalidate := models.AuthTokenSetter{
		ExpireAt: omitnull.From(time.Now()),
	}

	_, err := AuthTokens.Update(
		invalidate.UpdateMod(),
		um.Where(
			psql.And(
				AuthTokens.Columns.SessionID.EQ(psql.Arg(sessionID)),
				AuthTokens.Columns.Type.EQ(psql.Arg(tokenType)),
				AuthTokens.Columns.ExpireAt.GTE(psql.Arg(time.Now())),
			),
		),
	).All(ctx, r.db)

	if err != nil {
		return errors.Wrap(err, "error executing invalidate token query")
	}

	return nil
}

func NewAuthTokenRepository(db bob.Executor) *AuthTokenRepository {
	return &AuthTokenRepository{db: db}
}
