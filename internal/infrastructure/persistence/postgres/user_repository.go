package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"github.com/ztrue/tracerr"
)

var (
	_ repo.UserRepository = (*UserRepository)(nil)
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(
	db *pg.DB,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := u.db.ModelContext(ctx, user).Insert()

	if err != nil {
		if isMultiRowsError(err) || getConstraint(err) != "" {
			return nil, repo.ErrAlreadyExist
		}

		return nil, tracerr.Errorf("failed insert user: %w", err)
	}

	return user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := u.db.ModelContext(ctx, user).
		Where("id in (?)", user.ID).
		Update()

	if err != nil {
		if isNoRowsError(err) {
			return nil, repo.ErrNotFound
		}

		return nil, tracerr.Errorf("failed update user: %w", err)
	}

	return user, nil
}

func (u *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	result := &model.User{}
	q := u.db.ModelContext(ctx, result).
		Relation("Image").
		Where("\"user\".\"email\" = ?", email)

	if err := q.Select(); err != nil {
		if isNoRowsError(err) {
			return nil, repo.ErrNotFound
		}

		return nil, tracerr.Errorf("failed to find user: %w", err)
	}

	return result, nil
}

func (u *UserRepository) FindUserByID(ctx context.Context, id int) (*model.User, error) {
	result := &model.User{}
	q := u.db.ModelContext(ctx, result).
		Relation("Image").
		Where("\"user\".\"id\" = ?", id)

	if err := q.Select(); err != nil {
		if isNoRowsError(err) {
			return nil, repo.ErrNotFound
		}

		return nil, tracerr.Errorf("failed to find user: %w", err)
	}

	return result, nil
}

func (u *UserRepository) UploadImage(
	ctx context.Context,
	userId int,
	image, hash []byte,
	contentType string,
) error {
	newImage := model.Image{
		Hash:        hash,
		ContentType: contentType,
		File:        image,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	tx, err := u.db.Begin()
	if err != nil {
		tx.Rollback()
		return tracerr.Wrap(err)
	}

	_, err = tx.ModelContext(ctx, &newImage).Insert()
	if err != nil {
		tx.Rollback()
		return tracerr.Wrap(err)
	}

	_, err = tx.ModelContext(ctx, &model.User{}).
		Set("image_id = ?", newImage.ID).
		Where("id = ?", userId).
		Update()
	fmt.Println(err)

	if err != nil {
		tx.Rollback()
		return tracerr.Wrap(err)
	}

	tx.Commit()
	return nil
}
func (u *UserRepository) DownloadImage(ctx context.Context, userId int) ([]byte, string, error) {
	user := model.User{}

	err := u.db.ModelContext(ctx, &user).
		Relation("Image").
		Where("\"user\".\"id\" = ?", userId).
		Select()
	if err != nil {
		return nil, "", tracerr.Wrap(err)
	}
	return user.Image.File, user.Image.ContentType, nil
}
