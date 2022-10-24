package main

import (
	"context"
	"log"
	"security/internal/repositories/postgres"
	"security/internal/server"
	"security/internal/user"
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
