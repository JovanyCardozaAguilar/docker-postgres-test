package main

import (
	"context"
	"docker-go-test/data"
	"docker-go-test/models"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	ctx  context.Context
	pool *models.Postgres
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var middlewares = []Middleware{
	TokenAuthMiddleware,
}

func main() {
	dsn := "postgres://testUser1:password@localhost:5432/testdb1?sslmode=disable"
	ctx = context.Background()
	var err error
	pool, err = data.CreateDBPool(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("The Pool: ", pool)
	fmt.Println(data.QueryGreeting(ctx, pool))
	fmt.Println(data.QuerySingleTest(ctx, pool))
	accounts, _ := data.QueryMultiTest(ctx, pool)
	for _, account := range accounts {
		fmt.Printf("%#v\n", account)
	}

	var handler http.HandlerFunc = handleClientProfile
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	http.HandleFunc("/user/profile", handler)

	log.Println("Server is on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
