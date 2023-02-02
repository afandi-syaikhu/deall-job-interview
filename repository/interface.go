package repository

import (
	"context"

	"github.com/afandi-syaikhu/deall-job-interview/model"
)

//go:generate mockgen -destination=mock/user_mock.go -package=mock github.com/afandi-syaikhu/deall-job-interview/repository UserRepository
type UserRepository interface {
	FindByCredential(ctx context.Context, data model.User) (*model.User, error)
}
