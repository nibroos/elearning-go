package model

type User struct {
    ID              int32    `db:"id"`
    Name            string `db:"name"`
    Username        string `db:"username"`
    Email           string `db:"email"`
    Password        string `db:"password"`
}
