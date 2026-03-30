package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"

	"github.com/kubeploy/kubeploy/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthMiddleware struct {
	users  *models.UserStore
	cookie *securecookie.SecureCookie
}

func NewAuthMiddleware(users *models.UserStore, sessionSecret string) *AuthMiddleware {
	hashKey := []byte(sessionSecret)
	if len(hashKey) < 32 {
		hashKey = securecookie.GenerateRandomKey(32)
	}
	blockKey := securecookie.GenerateRandomKey(32)
	return &AuthMiddleware{
		users:  users,
		cookie: securecookie.New(hashKey, blockKey),
	}
}

func (m *AuthMiddleware) SetSession(w http.ResponseWriter, userID string) error {
	encoded, err := m.cookie.Encode("session", map[string]string{"user_id": userID})
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})
	return nil
}

func (m *AuthMiddleware) ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (m *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		values := map[string]string{}
		if err := m.cookie.Decode("session", cookie.Value, &values); err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		userID, ok := values["user_id"]
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		user, err := m.users.GetByID(userID)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(userContextKey).(*models.User)
	return user
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
