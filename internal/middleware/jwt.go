package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohdjishin/SplitWise/config"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
)

type contextKey string

const (
	ContextuserIdKey = contextKey("userId")
)

const authorization = "Authorization"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorization)
		if authHeader == "" {
			logger.LoggerInstance.Error("Authorization header not found")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrUnauthorizationHeaderNotFound)
			return
		}
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			logger.LoggerInstance.Error("Invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrInvalidAuthHeader)
			return
		}
		tokenString := strings.TrimSpace(parts[1])
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().JwtString), nil
		})

		if err != nil || !token.Valid {
			logger.LoggerInstance.Error("Invalid token", zap.Any("error", err))
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrInvalidToken)
			return
		}
		ctx := context.WithValue(r.Context(), ContextuserIdKey, (*claims)["id"].(float64))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCurrentUserId(r *http.Request) float64 {
	return r.Context().Value(ContextuserIdKey).(float64)
}
