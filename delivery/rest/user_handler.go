package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/afandi-syaikhu/deall-job-interview/constant"
	"github.com/afandi-syaikhu/deall-job-interview/delivery/middleware"
	"github.com/afandi-syaikhu/deall-job-interview/model"
	"github.com/afandi-syaikhu/deall-job-interview/usecase"
	"github.com/labstack/echo"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(e *echo.Echo, mw *middleware.Middleware, userUseCase usecase.UserUseCase) {
	handler := &UserHandler{
		UserUseCase: userUseCase,
	}

	// register route
	e.POST("v1/users", handler.Create, mw.MustAuth, mw.MustAdmin)
	e.GET("v1/users", handler.Fetch, mw.MustAuth)
	e.GET("v1/users/:id", handler.FindById, mw.MustAuth)
	e.PUT("v1/users/:id", handler.Update, mw.MustAuth, mw.MustAdmin)
	e.DELETE("v1/users/:id", handler.Delete, mw.MustAuth, mw.MustAdmin)
}

func (_u *UserHandler) Create(c echo.Context) error {
	body := model.CreateUserRequest{}
	response := model.Response{}
	ctx := c.Request().Context()

	err := c.Bind(&body)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	err = c.Validate(body)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	validRoles := _u.UserUseCase.GetValidRole()
	_, ok := validRoles[body.Role]
	if !ok {
		response.Message = constant.InvalidRole
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	err = _u.UserUseCase.Create(ctx, model.User{
		Username: body.Username,
		Password: body.Password,
		Role:     body.Role,
	})
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusInternalServerError, response)

		return echo.ErrInternalServerError
	}

	response.Message = constant.Success
	return c.JSON(http.StatusOK, response)
}

func (_u *UserHandler) Fetch(c echo.Context) error {
	response := model.Response{}
	ctx := c.Request().Context()

	paramPage := strings.TrimSpace(c.QueryParam(constant.KeyPage))
	if paramPage == "" {
		paramPage = fmt.Sprintf("%d", constant.DefaultPage)
	}

	paramLimit := strings.TrimSpace(c.QueryParam(constant.KeyLimit))
	if paramLimit == "" {
		paramLimit = fmt.Sprintf("%d", constant.DefaultLimit)
	}

	page, err := strconv.Atoi(paramPage)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return echo.ErrBadRequest
	}

	limit, err := strconv.Atoi(paramLimit)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return echo.ErrBadRequest
	}

	users, err := _u.UserUseCase.Fetch(ctx, model.PaginationRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusInternalServerError, response)

		return echo.ErrInternalServerError
	}

	usersResponse := []*model.UserResponse{}
	for _, user := range users {
		usersResponse = append(usersResponse, &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, usersResponse)
}

func (_u *UserHandler) FindById(c echo.Context) error {
	response := model.Response{}
	ctx := c.Request().Context()
	pathID := c.Param(constant.KeyId)

	if len(strings.TrimSpace(pathID)) == 0 {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	id, err := strconv.ParseInt(pathID, 10, 64)
	if err != nil {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	user, err := _u.UserUseCase.FindById(ctx, id)
	if err != nil && err.Error() == constant.NotFound {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	if err != nil {
		response.Message = constant.InternalServerError
		c.JSON(http.StatusInternalServerError, response)

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (_u *UserHandler) Update(c echo.Context) error {
	body := model.UpdateUserRequest{}
	response := model.Response{}
	ctx := c.Request().Context()
	pathID := c.Param(constant.KeyId)
	if len(strings.TrimSpace(pathID)) == 0 {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	id, err := strconv.ParseInt(pathID, 10, 64)
	if err != nil {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	err = c.Bind(&body)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	err = c.Validate(body)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	validRoles := _u.UserUseCase.GetValidRole()
	_, ok := validRoles[body.Role]
	if !ok {
		response.Message = constant.InvalidRole
		c.JSON(http.StatusBadRequest, response)

		return err
	}

	err = _u.UserUseCase.Update(ctx, model.User{
		ID:       id,
		Username: body.Username,
		Password: body.Password,
		Role:     body.Role,
	})
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusInternalServerError, response)

		return echo.ErrInternalServerError
	}

	response.Message = constant.Success
	return c.JSON(http.StatusOK, response)
}

func (_u *UserHandler) Delete(c echo.Context) error {
	response := model.Response{}
	ctx := c.Request().Context()
	pathID := c.Param(constant.KeyId)
	if len(strings.TrimSpace(pathID)) == 0 {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	id, err := strconv.ParseInt(pathID, 10, 64)
	if err != nil {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	requester, ok := c.Get(constant.KeyUser).(*model.User)
	if !ok || requester == nil {
		response.Message = constant.Unauthorized
		c.JSON(http.StatusUnauthorized, response)

		return echo.ErrUnauthorized
	}

	if requester.ID == id {
		response.Message = constant.NotAllowedSelfDelete
		c.JSON(http.StatusBadRequest, response)

		return echo.ErrBadRequest
	}

	err = _u.UserUseCase.DeleteById(ctx, id)
	if err != nil && err.Error() == constant.NotFound {
		response.Message = constant.NotFound
		c.JSON(http.StatusNotFound, response)

		return echo.ErrNotFound
	}

	if err != nil && err.Error() == constant.NotMeetMinimumAdmin {
		response.Message = constant.NotMeetMinimumAdmin
		c.JSON(http.StatusBadRequest, response)

		return echo.ErrBadRequest
	}

	if err != nil {
		response.Message = constant.InternalServerError
		c.JSON(http.StatusInternalServerError, response)

		return echo.ErrInternalServerError
	}

	response.Message = constant.Success
	return c.JSON(http.StatusOK, response)
}
