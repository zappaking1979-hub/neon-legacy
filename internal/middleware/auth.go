package middleware

import (
	"context"
	"net/http"

	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/domain/player"
)

type contextKey string

const (
	PlayerKey contextKey = "player"
)

func Auth(authService *application.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ""
			if c, err := r.Cookie("session"); err == nil {
				token = c.Value
			}

			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			p, err := authService.ValidateSession(r.Context(), token)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:   "session",
					Value:  "",
					Path:   "/",
					MaxAge: -1,
				})
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), PlayerKey, p)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetPlayer(r *http.Request) *player.Player {
	p, _ := r.Context().Value(PlayerKey).(*player.Player)
	return p
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := GetPlayer(r)
		if p == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireAuthAPI(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := GetPlayer(r)
		if p == nil {
			http.Error(w, `{"error":"not authenticated"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
