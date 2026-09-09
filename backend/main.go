package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/auth"
	"backend/config"
	"backend/handlers"
	"backend/middleware"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	sessionManager := auth.NewSessionManager()
	authHandler := handlers.NewAuthHandler(cfg, sessionManager)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./media"))
	mux.Handle("/media/", http.StripPrefix("/media/", fs))

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/verify", authHandler.Verify)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)

	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/verify", authHandler.Verify)
	mux.HandleFunc("/api/logout", authHandler.Logout)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("[!] http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.EnableCORS(mux)))
}
