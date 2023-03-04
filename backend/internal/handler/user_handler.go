package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"test-crud-user-orders/internal/entity"
	"test-crud-user-orders/internal/template"
	"time"

	"github.com/labstack/echo/v4"

	"test-crud-user-orders/internal/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

// Create Func for Inserting New Data
func (h *UserHandler) Create(c echo.Context) error {
	var input entity.CreateUser

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Bad Request",
		})
	}

	if err := c.Validate(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Bad Request",
		})
	}

	dataUser, err := h.userUseCase.Create(c.Request().Context(), input.FullName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, template.ResponseHTTP{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusCreated, template.ResponseHTTP{
		Status:  http.StatusCreated,
		Data:    dataUser,
		Message: "Create User Success",
	})
}

// GetAllPagination Func for Get All Data with Pagination func
func (h *UserHandler) GetAllPagination(c echo.Context) error {
	// init Pagination Logic
	limitData, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err != nil || limitData < 1 {
		limitData = 10 // default limit
	}
	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1 // default offset
	}
	// Calculate offset by limit per page and number of page
	offsetData := (page - 1) * limitData

	// Count Users Data, Return int64
	countData := h.userUseCase.CountData(c.Request().Context())
	var users []*entity.User
	if offsetData < countData {
		users, err = h.userUseCase.GetAllPagination(c.Request().Context(), int(limitData), int(offsetData))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, template.ResponseHTTP{
				Status:  http.StatusInternalServerError,
				Error:   err,
				Message: "Internal Server Error",
			})
		}
	}

	// Message for Result Data empty
	messageResult := "OK"
	if len(users) < 1 {
		messageResult = "Zero Data"
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Data:    users,
		Message: messageResult,
		Page: template.PagePagination{
			Limit: limitData,
			Page:  page,
			Show:  len(users),
			Total: countData,
		},
	})
}

// GetByID Func for Get 1 Data by primaryKey
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	user, err := h.userUseCase.GetByID(c.Request().Context(), int(id))
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("UserID %d Not Found or Deleted", id),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, template.ResponseHTTP{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: "Internal Server Error",
		})
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
			Status:  http.StatusNotFound,
			Data:    make([]int, 0, 0),
			Message: fmt.Sprintf("UserID %d Not Found or Deleted", id),
		})
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Data:    user,
		Message: "OK",
	})
}

// Update Func for Update 1 Data by primaryKey
func (h *UserHandler) Update(c echo.Context) error {
	// Condition IF client send unformatted PrimaryKey
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	// Condition IF client send Bad Request
	var input entity.CreateUser
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Bad Request",
		})
	}
	if err := c.Validate(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Bad Request",
		})
	}

	// Execute Update data of User by PrimaryKey
	if err := h.userUseCase.Update(c.Request().Context(), int(id), input.FullName); err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("UserID %d Not Found", id),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, template.ResponseHTTP{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("UserID %d Has Been Updated", id),
	})
}

// Delete Func for Delete 1 Data by primaryKey
func (h *UserHandler) Delete(c echo.Context) error {
	// Condition IF client send unformatted PrimaryKey
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	// Execute Delete data of User by PrimaryKey
	if err := h.userUseCase.Delete(c.Request().Context(), int(id)); err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("UserID %d Not Found", id),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, template.ResponseHTTP{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("UserID %d Has Been Deleted", id),
	})
}

func generateTime(days int) time.Time {
	today := time.Now()
	return today.AddDate(0, 0, days)
}
