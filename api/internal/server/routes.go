package server

import (
	"api/pkg/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func (s *Server) initRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/login", s.getLoginPage).Methods(http.MethodGet)
	r.HandleFunc("/login", s.login).Methods(http.MethodPost)

	r.HandleFunc("/signup", s.getSignupPage).Methods(http.MethodGet)
	r.HandleFunc("/signup", s.signup).Methods(http.MethodPost)

	r.HandleFunc("/logout", s.logout).Methods(http.MethodPost)

	CSRF := csrf.SimpleProtect(s.cfg.CSRFKey)
	r.Handle("/account", CSRF(http.HandlerFunc(s.getAccountPage))).Methods(http.MethodGet)
	// This is vulnerable to CSRF attack
	r.HandleFunc("/account/insecure", s.changeUsername).Methods(http.MethodPost)
	// This is protected from CSRF attack
	r.Handle("/account/secure", CSRF(http.HandlerFunc(s.changeUsername))).Methods(http.MethodPost)

	r.HandleFunc("/card", s.getCardInfo).Methods(http.MethodGet)

	configuredCORS := cors.New(
		cors.Options{
			AllowedOrigins: []string{"http://localhost:8080"},
			Debug:          true,
		},
	)
	//configuredCORS := cors.AllowAll()

	corsHandler := configuredCORS.Handler(r)
	return corsHandler
}
