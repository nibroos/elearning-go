package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"` // Added token field
}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=secret dbname=customerdb sslmode=disable")
	if err != nil {
		panic(err)
	}
}
