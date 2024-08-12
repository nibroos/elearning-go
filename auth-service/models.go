package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func init() {
    var err error
    db, err = sql.Open("postgres", "user=postgres password=secret dbname=authdb sslmode=disable")
    if err != nil {
        panic(err)
    }
}
