package main

import (
	"api/internal/repositories/postgres"
	"api/internal/server"
	"api/internal/user"
	"context"
	"log"
)

func main() {
	cfg := server.NewConfig()

	db, err := postgres.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUsers(db)

	userService := user.NewService(userRepo)

	s := server.NewServer(cfg, userService)

	log.Fatalf("Server stoped with error: %v", s.Start())
}
