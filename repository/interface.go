package repository

import (
	"context"

	"github.com/afandi-syaikhu/deall-job-interview/model"
)

//go:generate mockgen -destination=mock/user_mock.go -package=mock github.com/afandi-syaikhu/deall-job-interview/repository UserRepository
type UserRepository interface {
	Create(ctx context.Context, data model.User) error
	Fetch(ctx context.Context, data model.PaginationRequest) ([]*model.User, error)
	FindByCredential(ctx context.Context, data model.User) (*model.User, error)
	FindById(ctx context.Context, id int64) (*model.User, error)
	CountByRole(ctx context.Context, role model.Role) (int, error)
	IsExistByUsername(ctx context.Context, username string) (bool, error)
	IsExistByUsernameAndIdNot(ctx context.Context, username string, id int64) (bool, error)
	Update(ctx context.Context, data model.User) error
	DeleteById(ctx context.Context, id int64) error
}
