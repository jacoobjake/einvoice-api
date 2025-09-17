package repositories

import (
	"context"

	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/pkg/errors"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

var Users = models.Users

type UserRepository struct {
	db bob.Executor
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (*models.User, error) {
	return models.FindUser(ctx, r.db, id)
}

func (r *UserRepository) FindByIdOrFail(ctx context.Context, id int64) (*models.User, error) {
	user, err := r.FindById(ctx, id)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching user")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := Users.Query(
		sm.Where(Users.Columns.Email.EQ(psql.Arg(email))),
		sm.Where(Users.Columns.DeletedAt.IsNull()),
	).One(ctx, r.db)

	if err != nil {
		return nil, errors.Wrap(err, "error querying user")
	}

	return user, nil
}

func (r *UserRepository) FindByEmailOrFail(ctx context.Context, email string) (*models.User, error) {
	user, err := r.FindByEmail(ctx, email)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching user")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.UserSetter) (*models.User, error) {
	createdUser, err := Users.Insert(user).One(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting user record")
	}
	return createdUser, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User, data *models.UserSetter) (*models.User, error) {
	err := user.Update(ctx, r.db, data)
	if err != nil {
		return nil, errors.Wrap(err, "error updating user record")
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	return user.Delete(ctx, r.db)
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	users, err := Users.Query(
		sm.Limit(uint64(limit)),
		sm.Offset(uint64(offset)),
	).All(ctx, r.db)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching user list")
	}

	return users, nil
}

func NewUserRepository(db bob.Executor) *UserRepository {
	return &UserRepository{db: db}
}
