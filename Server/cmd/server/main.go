package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"mmorpg/internal/auth"
	"mmorpg/internal/config"
	"mmorpg/internal/game"
	"mmorpg/internal/gateway"
	"mmorpg/internal/storage"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
	})
}

// updatesHandler serves the update tree of the native client: manifests are always revalidated, objects are
// content-addressed and may be cached forever.
func updatesHandler(dir string) http.Handler {
	fs := http.StripPrefix("/updates/", http.FileServer(http.Dir(dir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/updates/objects/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fs.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()

	store, err := storage.OpenJSON(cfg.DataDir)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if len(cfg.Admins) > 0 {
		log.Printf("Admins (Map-Editor): %v", cfg.Admins)
	}
	world := game.New(cfg, store)
	worldDone := make(chan struct{})
	go func() { world.Run(ctx); close(worldDone) }()

	gw := gateway.New(store, auth.New(store), world, cfg.Admins)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", gw.ServeWS)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	if cfg.UpdatesDir != "" {
		mux.Handle("/updates/", updatesHandler(cfg.UpdatesDir))
		log.Printf("Client-Updates unter /updates/ aus %s", cfg.UpdatesDir)
	}
	if cfg.ClientDir != "" { // empty: API-only mode (client served by a separate container)
		mux.Handle("/", noCache(http.FileServer(http.Dir(cfg.ClientDir))))
	}

	srv := &http.Server{Addr: cfg.Addr, Handler: mux}
	go func() {
		abs, _ := filepath.Abs(cfg.ClientDir)
		log.Printf("Server läuft auf http://localhost%s (Client: %s)", cfg.Addr, abs)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Fahre herunter…")
	sctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = srv.Shutdown(sctx)
	<-worldDone
	if err := store.Close(); err != nil {
		log.Printf("store close: %v", err)
	}
}
