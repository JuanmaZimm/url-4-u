package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/JuanmaZimm/url-4-u/backend/internal/config"
	"github.com/JuanmaZimm/url-4-u/backend/internal/handler"
	"github.com/JuanmaZimm/url-4-u/backend/internal/repository"
	"github.com/JuanmaZimm/url-4-u/backend/internal/service"
)

func main() {
	cfg := config.Load()

	// Initialize SQLite Repository
	repo, err := repository.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database at %s: %v", cfg.DBPath, err)
	}
	defer repo.Close()
	log.Printf("📦 Connected to SQLite database: %s", cfg.DBPath)

	// Initialize Services
	urlService := service.NewURLService(repo, cfg.BaseURL)
	qrService := service.NewQRService()

	// Initialize Handlers
	urlHandler := handler.NewURLHandler(urlService, qrService, cfg.BaseURL)

	// Setup Router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"url-4-u","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
	})

	// Redirection Endpoint (GET and HEAD for link unfurlers)
	r.Get("/r/{code}", urlHandler.Redirect)
	r.Head("/r/{code}", urlHandler.Redirect)

	// API v1 Routes
	r.Route("/api/v1/urls", func(api chi.Router) {
		api.Post("/", urlHandler.Create)
		api.Get("/", urlHandler.List)
		api.Get("/{code}", urlHandler.Get)
		api.Get("/{code}/qr", urlHandler.GetQR)
		api.Get("/{code}/analytics", urlHandler.GetAnalytics)
	})

	// HTTP Server Configuration
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context for graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errorsIsTimeout := shutdownCtx.Err() == context.DeadlineExceeded; errorsIsTimeout {
				log.Fatal("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		log.Println("🛑 Shutting down server gracefully...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Printf("🚀 URL-4-U API server is running on http://localhost:%s", cfg.Port)
	log.Printf("🔗 Base URL configured as %s", cfg.BaseURL)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server startup failed: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Println("👋 Server gracefully exited")
}
