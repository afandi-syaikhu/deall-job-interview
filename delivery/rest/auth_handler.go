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
	response := model.Response{}
	body := model.LoginRequest{}
	err := c.Bind(&body)
	if err != nil {
		response.Message = constant.BadRequest
		c.JSON(http.StatusBadRequest, response)

		return echo.ErrBadRequest
	}

	if err := c.Validate(body); err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	user := model.User{
		Username: body.Username,
		Password: body.Password,
	}

	res, err := _a.AuthUseCase.Login(ctx, user)
	if err != nil && err.Error() == constant.InvalidCredential {
		response.Message = constant.InvalidCredential
		c.JSON(http.StatusUnauthorized, response)

		return echo.ErrUnauthorized
	}

	if err != nil {
		response.Message = constant.InternalServerError
		c.JSON(http.StatusInternalServerError, response)

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}
