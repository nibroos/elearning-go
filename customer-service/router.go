package main

import (
	"github.com/labstack/echo/v4"
)

func registerRoutes(e *echo.Echo) {
	e.POST("/customers/customers-index", getCustomersHandler)
	e.POST("/customers/customers-show", getCustomerHandler)
	e.POST("/customers/customers-store", registerCustomerHandler)
}
