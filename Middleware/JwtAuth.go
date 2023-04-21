package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const jwtSecret = "12ZEFRGHJK4RT5YUJIKIOLIuytreds"

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("no way")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if token is expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
			if time.Now().After(expirationTime) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
