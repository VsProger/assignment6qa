package cookie

import (
	"net/http"
	"time"
)

const cookieName = "session_id"

func GetSessionCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	return cookie
}

func SetSessionCookie(w http.ResponseWriter, token string, expirationTime time.Time) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		Expires:  expirationTime,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func ExpireSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}
