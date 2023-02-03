package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/afandi-syaikhu/deall-job-interview/constant"
	"github.com/afandi-syaikhu/deall-job-interview/model"
	"github.com/afandi-syaikhu/deall-job-interview/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

type Middleware struct {
	Config      *model.Config
	UserUseCase usecase.UserUseCase
}

func New(config *model.Config, userUC usecase.UserUseCase) *Middleware {
	return &Middleware{
		Config:      config,
		UserUseCase: userUC,
	}
}

func (_m *Middleware) MustAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			return httpError(http.StatusUnauthorized, constant.InvalidToken)
		}

		sAuthHeaders := strings.Split(authHeader, " ")
		if len(sAuthHeaders) != 2 || sAuthHeaders[0] != "Bearer" || sAuthHeaders[1] == "" {
			return httpError(http.StatusUnauthorized, constant.InvalidToken)
		}

		accessToken := sAuthHeaders[1]
		accessSecret := []byte(_m.Config.Jwt.AccessSecret)
		token, err := jwt.ParseWithClaims(accessToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(constant.ExpiredToken)
			}

			return accessSecret, nil
		})

		if err != nil {
			return httpError(http.StatusUnauthorized, err.Error())
		}

		claims, ok := token.Claims.(*model.Claims)
		if !ok || !token.Valid {
			return httpError(http.StatusUnauthorized, constant.InvalidToken)
		}

		user, err := _m.UserUseCase.FindById(c.Request().Context(), claims.UserId)
		if err != nil && err.Error() == constant.NotFound {
			return httpError(http.StatusUnauthorized, constant.UserNotExist)
		}

		if err != nil {
			return httpError(http.StatusInternalServerError, constant.InternalServerError)
		}

		c.Set(constant.KeyUser, user)

		return next(c)
	}
}

func (_m *Middleware) MustAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get(constant.KeyUser).(*model.User)
		if !ok || user == nil {
			return httpError(http.StatusUnauthorized, constant.Unauthorized)
		}

		if user.Role != constant.Admin {
			return httpError(http.StatusForbidden, constant.Forbidden)
		}

		return next(c)
	}
}

func (_m *Middleware) MustHaveRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			return httpError(http.StatusUnauthorized, constant.InvalidToken)
		}

		sAuthHeaders := strings.Split(authHeader, " ")
		if len(sAuthHeaders) != 2 || sAuthHeaders[0] != "Bearer" || sAuthHeaders[1] == "" {
			return httpError(http.StatusUnauthorized, constant.InvalidToken)
		}

		refreshToken := sAuthHeaders[1]
		refreshSecret := []byte(_m.Config.Jwt.RefreshSecret)
		token, err := jwt.ParseWithClaims(refreshToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(constant.ExpiredToken)
			}

			return refreshSecret, nil
		})

		if err != nil {
			return httpError(http.StatusUnauthorized, err.Error())
		}

		claims, ok := token.Claims.(*model.Claims)
		if !ok || !token.Valid {
			return httpError(http.StatusUnauthorized, constant.InvalidToken)
		}

		user, err := _m.UserUseCase.FindById(c.Request().Context(), claims.UserId)
		if err != nil && err.Error() == constant.NotFound {
			return httpError(http.StatusUnauthorized, constant.UserNotExist)
		}

		if err != nil {
			return httpError(http.StatusInternalServerError, constant.InternalServerError)
		}

		c.Set(constant.KeyUser, user)

		return next(c)
	}
}

func httpError(status int, message string) *echo.HTTPError {
	return echo.NewHTTPError(status, model.Response{
		Message: message,
	})
}
