package server

import "net/http"

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.cookieStore.Get(r, s.cfg.SessionName)

	// Revoke users authentication
	session.Values["authenticated"] = false
	delete(session.Values, "username")
	delete(session.Values, "login")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
