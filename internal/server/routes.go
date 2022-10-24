package server

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) initRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/login", s.getLoginPage).Methods(http.MethodGet)
	r.HandleFunc("/login", s.login).Methods(http.MethodPost)

	r.HandleFunc("/signup", s.getSignupPage).Methods(http.MethodGet)
	r.HandleFunc("/signup", s.signup).Methods(http.MethodPost)

	r.HandleFunc("/logout", s.logout).Methods(http.MethodPost)

	CSRF := csrf.Protect(s.cfg.CSRFKey, csrf.Secure(false))
	r.Handle("/account", CSRF(http.HandlerFunc(s.getAccountPage))).Methods(http.MethodGet)
	// This is vulnerable to CSRF attack
	r.HandleFunc("/account/insecure", s.changeUsername).Methods(http.MethodPost)
	// This is protected from CSRF attack
	r.Handle("/account/secure", CSRF(http.HandlerFunc(s.changeUsername))).Methods(http.MethodPost)

	return r
}
