package server

import (
	"api/internal/user"
	"net/http"
)

func (s *Server) getSignupPage(w http.ResponseWriter, r *http.Request) {
	if s.isAuthenticated(r) {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	s.render.HTML(w, http.StatusOK, "signup", nil)
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	session, _ := s.cookieStore.Get(r, s.cfg.SessionName)

	// Check if user is authenticated
	if s.isAuthenticated(r) {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	login := r.FormValue("login")
	username := s.nameGenerator.Generate()

	err := s.userService.Create(user.User{
		Username: username,
		Login:    login,
		Password: r.FormValue("password"),
	})
	if err != nil {
		s.render.String(w, http.StatusInternalServerError, err.Error())
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["login"] = login
	session.Save(r, w)
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	session, _ := s.cookieStore.Get(r, s.cfg.SessionName)
	if auth, ok := session.Values["authenticated"].(bool); auth && ok {
		return true
	}
	return false
}
