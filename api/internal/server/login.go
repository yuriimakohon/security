package server

import (
	"api/internal/user"
	"net/http"
)

func (s *Server) getLoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if s.isAuthenticated(r) {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}
	s.render.HTML(w, http.StatusOK, "login", nil)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	session, _ := s.cookieStore.Get(r, s.cfg.SessionName)

	login := r.FormValue("login")
	password := r.FormValue("password")

	usr, err := s.userService.Get(login)
	if err != nil {
		if err == user.ErrNotFound {
			s.render.String(w, http.StatusUnauthorized, "Invalid login or password")
			return
		}
		s.render.String(w, http.StatusInternalServerError, err.Error())
		return
	}

	if ok, err := s.userService.CheckPassword(login, password); err != nil {
		s.render.String(w, http.StatusInternalServerError, err.Error())
		return
	} else if !ok {
		s.render.String(w, http.StatusUnauthorized, "Invalid login or password")
		return
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["username"] = usr.Username
	session.Values["login"] = usr.Login
	session.Save(r, w)
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}
