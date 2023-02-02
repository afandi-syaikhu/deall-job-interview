package rest

import (
	"net/http"

	"github.com/afandi-syaikhu/deall-job-interview/constant"
	"github.com/afandi-syaikhu/deall-job-interview/model"
	"github.com/afandi-syaikhu/deall-job-interview/usecase"
	"github.com/labstack/echo"
)

type AuthHandler struct {
	AuthUseCase usecase.AuthUseCase
}

func NewAuthHandler(e *echo.Echo, authUseCase usecase.AuthUseCase) {
	handler := &AuthHandler{
		AuthUseCase: authUseCase,
	}

	// register route
	e.POST("/auth/login", handler.Login)
}

func (_a *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	errResponse := model.ErrorResponse{}
	body := model.LoginRequest{}
	if err := c.Bind(&body); err != nil {
		errResponse.Message = constant.BadRequest
		c.JSON(http.StatusBadRequest, errResponse)

		return echo.ErrBadRequest
	}

	if err := c.Validate(body); err != nil {
		errResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, errResponse)

		return err
	}

	user := model.User{
		Username: body.Username,
		Password: body.Password,
	}

	res, err := _a.AuthUseCase.Login(ctx, user)
	if err != nil && err.Error() == constant.InvalidCredential {
		errResponse.Message = constant.InvalidCredential
		c.JSON(http.StatusUnauthorized, errResponse)

		return echo.ErrUnauthorized
	}

	if err != nil {
		errResponse.Message = constant.InternalServerError
		c.JSON(http.StatusInternalServerError, errResponse)

		return echo.ErrInternalServerError
	}

	_a.setTokenCookie(res.AccessToken, c)
	_a.setUserCookie(&user, c)

	return c.JSON(http.StatusOK, res)
}

// Here we are creating a new cookie, which will store the valid JWT token.
func (_a *AuthHandler) setTokenCookie(token string, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "secret-key"
	cookie.Value = token
	// cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// Purpose of this cookie is to store the user's name.
func (_a *AuthHandler) setUserCookie(user *model.User, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Username
	// cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}
