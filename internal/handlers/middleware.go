package handlers

import (
	"forum/models"
	"forum/pkg/cookie"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

// func decorator(){

// }
func methodResolver(w http.ResponseWriter, r *http.Request, get, post func(w http.ResponseWriter, r *http.Request)) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	default:
		// error
	}
}

func (h *handler) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		c := cookie.GetSessionCookie(r)
		if c == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		isValid, err := h.service.ValidToken(c.Value)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		if !isValid {
			cookie.ExpireSessionCookie(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (h *handler) checkCookie(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cookie.GetSessionCookie(r)

		if c != nil {
			isValid, err := h.service.ValidToken(c.Value)
			if err != nil {
				h.app.ServerError(w, err)
				return
			}
			// TODO validate expire time of cookie

			if !isValid {
				cookie.ExpireSessionCookie(w)
				http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
				return
			}
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (h *handler) notRegistered(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cookie.GetSessionCookie(r)
		if c != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *handler) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func GetIntForm(r *http.Request, form string) (int, error) {
	valueString := r.FormValue(form)
	value, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (h *handler) NewTemplateData(r *http.Request) (*models.TemplateData, error) {
	var TemplateData models.TemplateData

	TemplateData.IsAuthenticated = h.isAuthenticated(r)

	if TemplateData.IsAuthenticated {
		user, err := h.service.GetUser(r)
		if err != nil {
			return nil, err
		}
		TemplateData.User = user
	}
	return &TemplateData, nil
}

func (h *handler) isAuthenticated(r *http.Request) bool {
	cookie := cookie.GetSessionCookie(r)
	return cookie != nil && cookie.Value != ""
}

func ConverCategories(CategoriesString []string) ([]int, error) {
	categories := make([]int, len(CategoriesString))
	for i, str := range CategoriesString {
		nb, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		categories[i] = nb
	}

	return categories, nil
}

func trim(s ...*string) {
	for _, value := range s {
		*value = strings.TrimSpace(*value)
	}
}
