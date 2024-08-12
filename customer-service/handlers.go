package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getCustomersHandler(c echo.Context) error {
	// Implement get customers logic
	// mock customers data
	customers := []Customer{
		{
			ID:    1,
			Name:  "John",
			Email: "bKuZi@example.com",
		},
		{
			ID:    2,
			Name:  "Jane",
			Email: "bKuZi@example.com",
		},
	}
	return c.JSON(http.StatusOK, customers)
}

func getCustomerHandler(c echo.Context) error {
	// Implement get customer by ID logic
	// mock customer data
	customer := Customer{
		ID:    1,
		Name:  "John",
		Email: "bKuZi@example.com",
	}

	return c.JSON(http.StatusOK, customer)
}

func registerCustomerHandler(c echo.Context) error {
	var customer Customer
	if err := c.Bind(&customer); err != nil {
		return err
	}

	// Generate a token for the customer
	token, err := generateToken()
	if err != nil {
		log.Println("Error generating token:", err)
		return err
	}
	customer.Token = token

	// Insert customer into the database with token
	go func() {
		_, err := db.Exec("INSERT INTO customers (name, email, token) VALUES ($1, $2, $3)", customer.Name, customer.Email, customer.Token)
		if err != nil {
			log.Println("Error inserting customer:", err)
		}
	}()

	return c.JSON(http.StatusCreated, customer)
}
