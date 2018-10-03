package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//AuthMiddleware : authentication middleware
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearer := strings.Split(authHeader, " ")
			if len(bearer) == 2 {
				token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Error")
					}
					return []byte("secret"), nil
				})
				if err != nil {
					json.NewEncoder(w).Encode(err.Error())
					return
				}
				if token.Valid {
					next(w, r)
				} else {
					resMap := make(map[string]string)
					resMap["message"] = "Failed"

					res, err := json.Marshal(resMap)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					w.WriteHeader(http.StatusBadRequest)
					w.Write(res)
				}
			} else {
				http.Error(w, "Invalid authorization header", http.StatusBadRequest)
				return
			}
		} else {
			resMap := make(map[string]string)
			resMap["message"] = "Failed"

			res, err := json.Marshal(resMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
		}
	})
}
