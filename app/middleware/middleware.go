package middleware

import (
	"app/pkg/models"
	"context"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

func SessionMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get("SessionID")
			if sessionID == "" {
				http.Error(w, "Unauthorized: Missing SessionID", http.StatusUnauthorized)
				return
			}

			var session models.Session
			if err := db.Where("id = ?", sessionID).First(&session).Error; err != nil {
				if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
					http.Error(w, "Unauthorized: Invalid SessionID", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if session.ExpiredAt.Before(time.Now()) {
				http.Error(w, "Unauthorized: Session Expired", http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, "userID", session.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
