package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"test-crud-user-orders/internal/entity"
	"test-crud-user-orders/internal/template"
	"test-crud-user-orders/internal/usecase"
)

type OrderItemHandler struct {
	orderItemUseCase usecase.OrderItemUseCase
}

func NewOrderItemHandler(orderItemUseCase usecase.OrderItemUseCase) *OrderItemHandler {
	return &OrderItemHandler{orderItemUseCase}
}

// Create Func for Inserting New Data
func (h *OrderItemHandler) Create(c echo.Context) error {
	var input entity.CreateOrderItem

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

	expiredAt, err := parseTime(input.ExpiredAt)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unkown Format expired_at",
		})
	}

	var orderItem entity.OrderItem
	orderItem.Name = input.Name
	orderItem.Price = input.Price
	orderItem.ExpiredAt = expiredAt

	if err = h.orderItemUseCase.Create(c.Request().Context(), &orderItem); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, template.ResponseHTTP{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusCreated, template.ResponseHTTP{
		Status:  http.StatusCreated,
		Message: "OK",
		Data:    orderItem,
	})
}

// GetAllPagination Func for Get All Data with Pagination func
func (h *OrderItemHandler) GetAllPagination(c echo.Context) error {
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
	countData := h.orderItemUseCase.CountData(c.Request().Context())
	var orderItem []*entity.OrderItem
	if offsetData < countData {
		orderItem, err = h.orderItemUseCase.GetAllPagination(c.Request().Context(), int(limitData), int(offsetData))
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
	lengthOrderItem := len(orderItem)
	if lengthOrderItem < 1 {
		messageResult = "Zero Data"
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Data:    orderItem,
		Message: messageResult,
		Page: template.PagePagination{
			Limit: limitData,
			Page:  page,
			Show:  lengthOrderItem,
			Total: countData,
		},
	})
}

// GetByID Func for Get 1 Data by primaryKey
func (h *OrderItemHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	orderItem, err := h.orderItemUseCase.GetByID(c.Request().Context(), int(id))
	if err != nil {
		fmt.Println("Error : " + err.Error())
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("OrderItemID #%d Not Found or Deleted", id),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, template.ResponseHTTP{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: "Internal Server Error",
		})
	}
	if orderItem == nil {
		return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
			Status:  http.StatusNotFound,
			Data:    make([]int, 0, 0),
			Message: fmt.Sprintf("OrderItemID #%d Not Found or Deleted", id),
		})
	}

	return c.JSON(http.StatusOK, template.ResponseHTTP{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    orderItem,
	})
}

// Update Func for Update 1 Data by primaryKey
func (h *OrderItemHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	var input entity.CreateOrderItem

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
			Message: "Error Validate Request",
		})
	}

	expiredAt, err := parseTime(input.ExpiredAt)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Error Parse Request expired_at",
		})
	}

	var orderItem entity.OrderItem
	orderItem.ID = id
	orderItem.Name = input.Name
	orderItem.Price = input.Price
	orderItem.ExpiredAt = expiredAt

	if err := h.orderItemUseCase.Update(c.Request().Context(), &orderItem); err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("OrderItemID #%d Not Found or Deleted", id),
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
		Message: "OK",
		Data:    orderItem,
	})
}

// Delete Func for Delete 1 Data by primaryKey
func (h *OrderItemHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, template.ResponseHTTP{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: "Unknown ID",
		})
	}

	if err := h.orderItemUseCase.Delete(c.Request().Context(), int(id)); err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, template.ResponseHTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("OrderItemID #%d Not Found", id),
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
		Message: fmt.Sprintf("OrderItemID #%d Has Been Deleted", id),
	})
}
