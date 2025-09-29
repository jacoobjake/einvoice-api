package repositories

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/jacoobjake/einvoice-api/pkg"
	"github.com/pkg/errors"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

var FailedLogins = models.FailedLogins

type FailedLoginRepository struct {
	db bob.Executor
}

func (r *FailedLoginRepository) Create(ctx context.Context, fl *models.FailedLoginSetter) (*models.FailedLogin, error) {
	failedLogin, err := FailedLogins.Insert(fl).One(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting failed_logins")
	}
	return failedLogin, nil
}

func (r *FailedLoginRepository) CaptureFailedLogin(ctx context.Context, userId int64) (*models.FailedLogin, error) {
	clientIp, ok := pkg.GetCtxClientIp(ctx)

	if !ok {
		return nil, errors.New("invalid client ip")
	}
	data := &models.FailedLoginSetter{
		UserID:    omit.From(userId),
		IPAddress: omit.From(clientIp),
	}

	return r.Create(ctx, data)
}

func (r *FailedLoginRepository) GetFailedLoginCountByUserId(ctx context.Context, userId int64) (int64, error) {
	failedLogins, err := FailedLogins.Query(
		sm.Where(FailedLogins.Columns.UserID.EQ(psql.Arg(userId))),
	).Count(ctx, r.db)

	if err != nil {
		return 0, errors.New("error querying failed_logins count by userId")
	}

	return failedLogins, nil
}

func NewFailedLoginRepository(db bob.Executor) *FailedLoginRepository {
	return &FailedLoginRepository{
		db: db,
	}
}
