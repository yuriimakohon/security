package server

import (
	"api/internal/user"
	"github.com/goombaio/namegenerator"
	"github.com/gorilla/sessions"
	"github.com/thedevsaddam/renderer"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer http.Server
	cfg        Config

	cookieStore *sessions.CookieStore
	render      *renderer.Render

	nameGenerator namegenerator.Generator

	userService *user.Service
}

func NewServer(cfg Config, user *user.Service) *Server {
	return &Server{
		httpServer: http.Server{
			Addr: ":" + cfg.Port,
		},
		cfg:         cfg,
		cookieStore: sessions.NewCookieStore(cfg.SessionKey),
		render: renderer.New(renderer.Options{
			ParseGlobPattern: cfg.TemplatesDir + "/*.html",
		}),
		nameGenerator: namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()),
		userService:   user,
	}
}

func (s *Server) Start() error {
	log.Println("Starting server on port", s.cfg.Port)

	s.httpServer.Handler = s.initRoutes()

	return s.httpServer.ListenAndServe()
}
