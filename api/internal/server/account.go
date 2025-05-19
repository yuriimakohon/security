package server

import (
	"api/pkg/csrf"
	"net/http"
)

func (s *Server) getAccountPage(w http.ResponseWriter, r *http.Request) {
	session, _ := s.cookieStore.Get(r, s.cfg.SessionName)

	if !s.isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Can't get name", http.StatusInternalServerError)
		return
	}

	s.render.HTML(w, http.StatusOK, "account", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"username":       username,
	})
}

func (s *Server) changeUsername(w http.ResponseWriter, r *http.Request) {
	session, _ := s.cookieStore.Get(r, s.cfg.SessionName)

	if !s.isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Empty username", http.StatusBadRequest)
		return
	}

	login, ok := session.Values["login"].(string)
	if !ok {
		http.Error(w, "Can't get login", http.StatusInternalServerError)
		return
	}

	err := s.userService.UpdateUsername(login, username)
	if err != nil {
		s.render.String(w, http.StatusInternalServerError, err.Error())
		return
	}

	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (s *Server) getCardInfo(w http.ResponseWriter, r *http.Request) {
	//if !s.isAuthenticated(r) {
	//	http.Redirect(w, r, "/login", http.StatusSeeOther)
	//	return
	//}

	cardInfo, err := s.userService.GetCardInfo()
	if err != nil {
		s.render.String(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"cardInfo": "` + cardInfo + `"}`))
}
