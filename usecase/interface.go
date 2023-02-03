package usecase

import (
	"context"

	"github.com/afandi-syaikhu/deall-job-interview/model"
)

//go:generate mockgen -destination=mock/auth_mock.go -package=mock github.com/afandi-syaikhu/deall-job-interview/usecase AuthUseCase
type AuthUseCase interface {
	Login(ctx context.Context, data model.User) (*model.Token, error)
	RefreshToken(ctx context.Context, data model.User) (*model.Token, error)
}

//go:generate mockgen -destination=mock/user_mock.go -package=mock github.com/afandi-syaikhu/deall-job-interview/usecase UserUseCase
type UserUseCase interface {
	Create(ctx context.Context, data model.User) error
	Fetch(ctx context.Context, data model.PaginationRequest) ([]*model.User, error)
	FindById(ctx context.Context, id int64) (*model.User, error)
	GetValidRole() map[model.Role]model.Role
	Update(ctx context.Context, data model.User) error
	DeleteById(ctx context.Context, id int64) error
}
