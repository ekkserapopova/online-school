package middleware

import (
	"context"
	"go.uber.org/fx"
	"net/http"
	"onlineschool/internal/pkg/jwter"
	"onlineschool/pkg/responser"
	"strings"
)

type AuthMiddlewareParams struct {
	fx.In

	JWTer *jwter.JWTer
}

type AuthMiddleware struct {
	jwt *jwter.JWTer
}

func NewAuthMiddleware(p AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: p.JWTer,
	}
}

func (authMD *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			responser.SendErr(w, http.StatusForbidden, "not authorized")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			responser.SendErr(w, http.StatusForbidden, "Authorization header format must be Bearer {token}")
			return
		}

		token := parts[1]

		claims, err := authMD.jwt.ValidateJWT(token)
		if err != nil {
			responser.SendErr(w, http.StatusForbidden, "Invalid token")
			return
		}

		// Получаем ID пользователя из claims
		userID, err := authMD.jwt.GetUserID(claims)
		if err != nil {
			responser.SendErr(w, http.StatusForbidden, "Invalid user ID in token")
			return
		}

		// Добавляем userID в контекст с простым строковым ключом
		ctx := context.WithValue(r.Context(), "userID", userID)

		// Продолжаем выполнение запроса с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
