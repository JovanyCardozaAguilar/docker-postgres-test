package main

import (
	"context"
	"docker-go-test/data"
	"docker-go-test/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	ctx  context.Context
	pool *models.Postgres
)

func handleClientProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HANDLE CLIENT SPOT: ", r.Method)
	switch r.Method {
	case http.MethodGet:
		GetClientProfile(w, r)
	case http.MethodPatch:
		UpdateClientProfile(w, r)
	case http.MethodPut:
		PutClientProfile(w, r)
	case http.MethodDelete:
		DeleteClientProfile(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value("clientProfile").(*models.ClientProfile)

	w.Header().Set("Content-Type", "application/json")

	response := models.ClientProfile{
		Id:        clientProfile.Id,
		FirstName: clientProfile.FirstName,
		LastName:  clientProfile.LastName,
		Token:     clientProfile.Token,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value("clientProfile").(*models.ClientProfile)

	// Decode the JSON payload into struct
	var payloadData models.ClientProfile
	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	clientProfile.FirstName = payloadData.FirstName
	clientProfile.LastName = payloadData.LastName
	clientProfile.Token = payloadData.Token
	fmt.Println("The payload data: ", payloadData)
	fmt.Println("The changed Client Profile: ", clientProfile)
	data.UpdateUser(pool, ctx, clientProfile.Id, *clientProfile)

	w.WriteHeader(http.StatusOK)
}

func PutClientProfile(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload into struct
	var payloadData models.ClientProfile
	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("The payload data: ", payloadData)
	data.InsertUser(pool, ctx, payloadData)

	w.WriteHeader(http.StatusOK)
}

func DeleteClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value("clientProfile").(*models.ClientProfile)
	data.DeleteUser(pool, ctx, clientProfile.Id)

	w.WriteHeader(http.StatusOK)
}

func TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("clientId")
		id, _ := strconv.Atoi(clientId)
		if id == 0 {
			fmt.Println("The Id", id)
			next.ServeHTTP(w, r)
			return
		}
		clientProfile, ok := data.GetUser(pool, ctx, id)
		if ok != nil || clientId == "" {
			http.Error(w, "ClientID does not exist Forbidden", http.StatusForbidden)
			return
		}
		token := r.Header.Get("Authorization")
		if !isValidToken(*clientProfile, token) {
			http.Error(w, "Authorization Forbidden", http.StatusForbidden)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "clientProfile", clientProfile))
		next.ServeHTTP(w, r)
	}
}

func isValidToken(clientProfile models.ClientProfile, token string) bool {
	if strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer ") == clientProfile.Token
	}
	return false
}

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
