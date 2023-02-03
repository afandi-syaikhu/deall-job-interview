package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/afandi-syaikhu/deall-job-interview/constant"
	"github.com/afandi-syaikhu/deall-job-interview/model"
	"github.com/afandi-syaikhu/deall-job-interview/repository"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type Auth struct {
	UserRepo repository.UserRepository
	Config   *model.Config
}

func NewAuthUseCase(userRepo repository.UserRepository, config *model.Config) AuthUseCase {
	return &Auth{
		UserRepo: userRepo,
		Config:   config,
	}
}

func (_a *Auth) Login(ctx context.Context, data model.User) (*model.Token, error) {
	user, err := _a.UserRepo.FindByCredential(ctx, data)
	if err != nil && err == sql.ErrNoRows {
		return nil, errors.New(constant.InvalidCredential)
	}

	if err != nil {
		log.Errorf("[%s] => %s", "AuthUC.Login", err.Error())
		return nil, err
	}

	accessToken, err := _a.generateJWT(*user, _a.Config.Jwt.AccessSecret, _a.Config.Jwt.AccessExp)
	if err != nil {
		log.Errorf("[%s] => %s", "AuthUC.Login", err.Error())
		return nil, err
	}

	refreshToken, err := _a.generateJWT(*user, _a.Config.Jwt.RefreshSecret, _a.Config.Jwt.RefreshExp)
	if err != nil {
		log.Errorf("[%s] => %s", "AuthUC.Login", err.Error())
		return nil, err
	}

	return &model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         constant.BearerType,
	}, nil
}

func (_a *Auth) generateJWT(user model.User, secretKey string, expTime int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expTime)).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
