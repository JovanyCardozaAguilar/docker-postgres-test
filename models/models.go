package models

import "github.com/jackc/pgx/v5/pgxpool"

type Postgres struct {
	Db *pgxpool.Pool
}

type ClientProfile struct {
	Id        int
	FirstName string
	LastName  string
	Token     string
}
