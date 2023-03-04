package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"test-crud-user-orders/internal/entity"
	"test-crud-user-orders/internal/template"

	"github.com/labstack/echo/v4"

	"test-crud-user-orders/internal/usecase"
)

type OrderHistoryHandler struct {
	orderHistoryUseCase usecase.OrderHistoryUseCase
}

func NewOrderHistoryHandler(orderHistoryUseCase usecase.OrderHistoryUseCase) *OrderHistoryHandler {
	return &OrderHistoryHandler{orderHistoryUseCase}
}

// Create Func for Inserting New Data
func (h *OrderHistoryHandler) Create(c echo.Context) error {
	var orderHistory *entity.OrderHistory
	var input entity.CreateOrderHistory
	var err error

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Invalid Request",
		})
	}

	if err := c.Validate(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Bad Request",
		})
	}

	if orderHistory, err = h.orderHistoryUseCase.Create(c.Request().Context(), input.UserID, input.OrderItemID, input.Descriptions); err != nil {
		var errResult error
		message := "Internal Server Error"
		status := http.StatusInternalServerError

		if err.Error() == "user not found" {
			status = http.StatusNotFound
			message = fmt.Sprintf("UserID %d Not Found or Deleted", input.UserID)
		} else if err.Error() == "order item not found" {
			status = http.StatusNotFound
			message = fmt.Sprintf("OrderItemID %d Not Found or Deleted", input.OrderItemID)
		} else {
			errResult = err
		}

		return echo.NewHTTPError(status, template.ResponseHTTP{
			Status:  status,
			Error:   errResult,
			Message: message,
		})
	}

	return c.JSON(http.StatusCreated, template.ResponseHTTP{
		Status:  http.StatusCreated,
		Data:    orderHistory,
		Message: "OK",
	})
}

// GetAllPagination Func for Get All Data with Pagination func
func (h *OrderHistoryHandler) GetAllPagination(c echo.Context) error {
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
	countData := h.orderHistoryUseCase.CountData(c.Request().Context(), -1)
	var orderHistory []*entity.OrderHistory
	if offsetData < countData {
		orderHistory, err = h.orderHistoryUseCase.GetAllPagination(c.Request().Context(), int(limitData), int(offsetData))
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
	lenOrderHistory := len(orderHistory)
	if lenOrderHistory < 1 {
		messageResult = "Zero Data"
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Data:    orderHistory,
		Message: messageResult,
		Page: template.PagePagination{
			Limit: limitData,
			Page:  page,
			Show:  lenOrderHistory,
			Total: countData,
		},
	})
}

// GetByID Func for Get 1 Data by primaryKey
func (h *OrderHistoryHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	orderHistory, err := h.orderHistoryUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		fmt.Println("Error : " + err.Error())
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: "Order History Not Found",
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
		Data:    orderHistory,
		Message: "OK",
	})
}

// Update Func for Update 1 Data by primaryKey
func (h *OrderHistoryHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	var input entity.CreateOrderHistory

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Invalid Request",
		})
	}

	if err := c.Validate(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Bad Request",
		})
	}

	if err := h.orderHistoryUseCase.Update(c.Request().Context(), id, input.UserID, input.OrderItemID, input.Descriptions); err != nil {
		var errDB error
		status := http.StatusInternalServerError
		message := "Internal Server Error"

		if err.Error() == "record not found" || err.Error() == "order history not found" {
			status = http.StatusNotFound
			message = "Order History Not Found"
		} else if err.Error() == "user data not found" {
			status = http.StatusNotFound
			message = "UserID Not Found"
		} else if err.Error() == "order item data not found" {
			status = http.StatusNotFound
			message = "OrderItemID Not Found"
		} else {
			errDB = err
		}

		return echo.NewHTTPError(status, template.ResponseHTTP{
			Status:  status,
			Error:   errDB,
			Message: message,
		})
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusBadRequest,
		Message: "OK",
	})
}

// Delete Func for Delete 1 Data by primaryKey
func (h *OrderHistoryHandler) Delete(c echo.Context) error {
	return echo.NewHTTPError(http.StatusForbidden, template.ResponseHTTP{
		Status:  http.StatusForbidden,
		Message: "Delete Transaction Not Allowed",
	})
}

// GetHistoryByUserID Func for Get All History Data by UserID with Pagination func
func (h *OrderHistoryHandler) GetHistoryByUserID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

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
	countData := h.orderHistoryUseCase.CountData(c.Request().Context(), userID)
	var orderHistory []*entity.OrderHistory
	if offsetData < countData {
		orderHistory, err = h.orderHistoryUseCase.GetByUserID(c.Request().Context(), userID, int(limitData), int(offsetData))
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
	lenOrderHistory := len(orderHistory)
	if lenOrderHistory < 1 {
		messageResult = "Zero Data"
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Data:    orderHistory,
		Message: messageResult,
		Page: template.PagePagination{
			Limit: limitData,
			Page:  page,
			Show:  lenOrderHistory,
			Total: countData,
		},
	})
}
