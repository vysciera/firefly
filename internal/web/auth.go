package web

import (
	"errors"
	"net/http"

	"firefly/internal/flowerpress"
)

type loginPageData struct {
	Error string
}

func sessionFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(flowerpress.SessionCookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func relayCookies(w http.ResponseWriter, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(
		w, "login.html",
		loginPageData{},
	); err != nil {
		http.Error(
			w, "failed to render login",
			http.StatusInternalServerError,
		)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(
			w, "invalid form",
			http.StatusBadRequest,
		)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	_, cookies, err := s.flowerpress.Login(
		r.Context(),
		username,
		password,
	)

	relayCookies(w, cookies)
	if errors.Is(err, flowerpress.ErrInvalidCredentials) {
		w.WriteHeader(http.StatusUnauthorized)

		_ = s.templates.ExecuteTemplate(
			w, "login.html",
			loginPageData{
				Error: "invalid credentials",
			},
		)
		return
	}

	if err != nil {
		http.Error(
			w, "flowerpress unavailable",
			http.StatusBadGateway,
		)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	session := sessionFromRequest(r)
	cookies, err := s.flowerpress.Logout(r.Context(), session)

	relayCookies(w, cookies)
	if err != nil {
		http.Error(
			w, "logout failed",
			http.StatusBadGateway,
		)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
