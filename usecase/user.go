package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/afandi-syaikhu/deall-job-interview/constant"
	"github.com/afandi-syaikhu/deall-job-interview/model"
	"github.com/afandi-syaikhu/deall-job-interview/repository"
	log "github.com/sirupsen/logrus"
)

type User struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &User{
		UserRepo: userRepo,
	}
}

func (_a *User) Create(ctx context.Context, data model.User) error {
	isExist, err := _a.UserRepo.IsExistByUsername(ctx, data.Username)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.Create", err.Error())
		return err
	}

	if isExist {
		return errors.New(constant.UsernameExist)
	}

	err = _a.UserRepo.Create(ctx, data)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.Create", err.Error())
		return err
	}

	return nil
}

func (_a *User) Fetch(ctx context.Context, data model.PaginationRequest) ([]*model.User, error) {
	if data.Limit <= 0 || data.Page <= 0 {
		return nil, errors.New(constant.BadRequest)
	}

	users, err := _a.UserRepo.Fetch(ctx, data)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.Fetch", err.Error())
		return nil, err
	}

	return users, nil
}

func (_a *User) FindById(ctx context.Context, id int64) (*model.User, error) {
	user, err := _a.UserRepo.FindById(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		return nil, errors.New(constant.NotFound)
	}

	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.FindById", err.Error())
		return nil, err
	}

	return user, nil
}

func (_a *User) GetValidRole() map[model.Role]model.Role {
	return map[model.Role]model.Role{
		constant.Admin: constant.Admin,
		constant.User:  constant.User,
	}
}

func (_a *User) Update(ctx context.Context, data model.User) error {
	isExist, err := _a.UserRepo.IsExistByUsernameAndIdNot(ctx, data.Username, data.ID)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.Update", err.Error())
		return err
	}

	if isExist {
		return errors.New(constant.UsernameExist)
	}

	user, err := _a.FindById(ctx, data.ID)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.Update", err.Error())
		return err
	}

	if user.Role == constant.Admin && data.Role != constant.Admin {
		totalAdmin, err := _a.UserRepo.CountByRole(ctx, constant.Admin)
		if err != nil {
			log.Errorf("[%s] => %s", "UserUC.Update", err.Error())
			return err
		}

		if totalAdmin < constant.MinimumAdmin+1 {
			return errors.New(constant.NotMeetMinimumAdmin)
		}
	}

	err = _a.UserRepo.Update(ctx, data)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.Update", err.Error())
		return err
	}

	return nil
}

func (_a *User) DeleteById(ctx context.Context, id int64) error {
	user, err := _a.FindById(ctx, id)
	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.DeleteById", err.Error())
		return err
	}

	if user.Role == constant.Admin {
		totalAdmin, err := _a.UserRepo.CountByRole(ctx, constant.Admin)
		if err != nil {
			log.Errorf("[%s] => %s", "UserUC.DeleteById", err.Error())
			return err
		}

		if totalAdmin < constant.MinimumAdmin+1 {
			return errors.New(constant.NotMeetMinimumAdmin)
		}
	}

	err = _a.UserRepo.DeleteById(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		return errors.New(constant.NotFound)
	}

	if err != nil {
		log.Errorf("[%s] => %s", "UserUC.DeleteById", err.Error())
		return err
	}

	return nil
}
