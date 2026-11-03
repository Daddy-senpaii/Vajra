package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Daddy-senpaii/Vajra/services/api/internal/cleanbrowsing"
	"github.com/Daddy-senpaii/Vajra/services/api/internal/config"
	"github.com/Daddy-senpaii/Vajra/services/api/internal/httpapi"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	domains, err := cleanbrowsing.CuratedDomains()
	if err != nil {
		log.Fatalf("load curated domains: %v", err)
	}
	provider := cleanbrowsing.NewClient(cfg.CleanBrowsingAPIURL, cfg.CleanBrowsingAPIKey, cfg.CleanBrowsingProfile, nil)
	var providerClient *cleanbrowsing.Client
	if cfg.CleanBrowsingAPIKey != "" {
		providerClient = provider
	}

	api := httpapi.New(httpapi.Options{
		Domains: domains,
		Config: httpapi.Config{
			DoTHost:  cfg.CleanBrowsingDoTHost,
			DoHURL:   cfg.CleanBrowsingDoHURL,
			Filter:   "adult",
			Provider: "CleanBrowsing",
			Profile:  cfg.CleanBrowsingProfile,
		},
		AdminToken:     cfg.AdminToken,
		AllowedOrigins: cfg.AllowedOrigins,
		CleanBrowsing:  providerClient,
	})
	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           api.Handler(),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	}()

	log.Printf("Vajra API listening on %s with %d curated domains", cfg.ListenAddr, len(domains))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
