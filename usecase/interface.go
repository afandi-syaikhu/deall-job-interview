package usecase

import (
	"context"

	"github.com/afandi-syaikhu/deall-job-interview/model"
)

//go:generate mockgen -destination=mock/auth_mock.go -package=mock github.com/afandi-syaikhu/deall-job-interview/usecase AuthUseCase
type AuthUseCase interface {
	Login(ctx context.Context, data model.User) (*model.Token, error)
}
