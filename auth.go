package main

import (
	"context"
	"docker-go-test/data"
	"docker-go-test/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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
