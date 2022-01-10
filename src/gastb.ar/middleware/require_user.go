package middleware

// The middleware package checks if a user is logged in and reacts accordingly.

import (
	"net/http"

	"gastb.ar/models"
	"gastb.ar/context"
)

// RequireUser wraps the UserService and adds verification methods
type RequireUser struct {
	*models.UserService
}

// ApplyFn takes in a handler function and returns it again only if user
// is logged in; otherwise, it redirects to login page
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		
		user, err := mw.UserService.ByToken(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next(w, r)
	})
}

// Apply takes in a handler and passes its ServeHTTP handler function
// over to RequireFn
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
