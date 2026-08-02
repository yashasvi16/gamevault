package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"encoding/json"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" || !strings.HasPrefix(authorization, "Bearer"){
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Authorization token is required",
			})
			return
		}
		
		tokenString := strings.TrimPrefix(authorization, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid authorization token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid authorization token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "player_id", claims["player_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}