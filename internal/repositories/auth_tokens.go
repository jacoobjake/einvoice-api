package repositories

import (
	"context"
	"time"

	"github.com/jacoobjake/einvoice-api/internal/database/enums"
	"github.com/jacoobjake/einvoice-api/internal/database/models"
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
		return nil, err
	}
	return createdToken, nil
}

func (r *AuthTokenRepository) FindByToken(ctx context.Context, token string) (*models.AuthToken, error) {
	authToken, err := AuthTokens.Query(
		sm.Where(AuthTokens.Columns.Token.EQ(psql.Arg(token))),
	).One(ctx, r.db)

	if err != nil {
		return nil, err
	}

	return authToken, nil
}

func (r *AuthTokenRepository) FindTokenByUserIdAndType(ctx context.Context, userID int64, tokenType enums.AuthTokenTypes) (*models.AuthToken, error) {
	authToken, err := AuthTokens.Query(
		sm.Where(AuthTokens.Columns.UserID.EQ(psql.Arg(userID))),
		sm.Where(AuthTokens.Columns.Type.EQ(psql.Arg(tokenType))),
	).One(ctx, r.db)

	if err != nil {
		return nil, err
	}

	return authToken, nil
}

func (r *AuthTokenRepository) InvalidateActiveTokensByUserID(ctx context.Context, userID int64, tokenType enums.AuthTokenTypes) error {
	_, err := AuthTokens.Update(
		um.Set(AuthTokens.Columns.ExpireAt, psql.Arg(time.Now())),
		um.Where(
			AuthTokens.Columns.UserID.EQ(psql.Arg(userID)),
		),
		um.Where(AuthTokens.Columns.Type.EQ(psql.Arg(tokenType))),
		um.Where(AuthTokens.Columns.ExpireAt.GT(psql.Arg(time.Now()))),
	).All(ctx, r.db)

	if err != nil {
		return err
	}

	return nil
}
