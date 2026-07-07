package main

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/neonlegacy/server/internal/adapters/postgres"
	"github.com/neonlegacy/server/internal/adapters/redis"
	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/config"
	"github.com/neonlegacy/server/internal/handlers"
	"github.com/neonlegacy/server/internal/middleware"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := postgres.NewPool(ctx, cfg.DB.URL, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.ConnMaxLifetime)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer dbPool.Close()
	log.Println("connected to postgres")

	rdb, err := redis.NewClient(ctx, cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer rdb.Close()
	log.Println("connected to redis")

	if err := runMigrations(cfg.DB.URL); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations up to date")

	playerRepo := postgres.NewPlayerRepo(dbPool)
	authService := application.NewAuthService(playerRepo, rdb, cfg.Auth.SessionTTL, cfg.Auth.BcryptCost)

	crimeRepo := postgres.NewCrimeRepo(dbPool)
	crimeService := application.NewCrimeService(crimeRepo, playerRepo)

	gymRepo := postgres.NewGymRepo(dbPool)
	gymService := application.NewGymService(gymRepo, playerRepo)

	itemRepo := postgres.NewItemRepo(dbPool)
	playerItemRepo := postgres.NewPlayerItemRepo(dbPool)
	shopService := application.NewShopService(itemRepo, playerItemRepo, playerRepo)
	inventoryService := application.NewInventoryService(playerItemRepo, playerRepo)

	tmpl := loadTemplates("web/templates")

	authHandler := handlers.NewAuthHandler(authService, tmpl)
	crimeHandler := handlers.NewCrimeHandler(crimeService, playerRepo, tmpl)
	gymHandler := handlers.NewGymHandler(gymService, playerRepo, tmpl)
	shopHandler := handlers.NewShopHandler(shopService, playerRepo, tmpl)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, playerRepo, tmpl)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.StripSlashes)
	r.Use(httprate.LimitByIP(cfg.Auth.RateLimit, cfg.Auth.RateLimitWindow))
	r.Use(middleware.Auth(authService))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", authHandler.LoginPage)
		r.Post("/login", authHandler.Login)
		r.Get("/register", authHandler.RegisterPage)
		r.Post("/register", authHandler.Register)
		r.Post("/logout", authHandler.Logout)
	})

	r.Get("/login", authHandler.LoginPage)
	r.Get("/register", authHandler.RegisterPage)
	r.Post("/login", authHandler.Login)
	r.Post("/register", authHandler.Register)
	r.Post("/logout", authHandler.Logout)
	r.Get("/dashboard", middleware.RequireAuth(authHandler.Dashboard))
	r.Get("/crimes", middleware.RequireAuth(crimeHandler.Page))
	r.Post("/crimes/commit", middleware.RequireAuth(crimeHandler.DoCrime))
	r.Get("/gym", middleware.RequireAuth(gymHandler.Page))
	r.Post("/gym/train", middleware.RequireAuth(gymHandler.Train))
	r.Get("/shop", middleware.RequireAuth(shopHandler.Page))
	r.Post("/shop/buy", middleware.RequireAuth(shopHandler.Buy))
	r.Post("/shop/sell", middleware.RequireAuth(shopHandler.Sell))
	r.Get("/inventory", middleware.RequireAuth(inventoryHandler.Page))
	r.Post("/inventory/use", middleware.RequireAuth(inventoryHandler.Use))
	r.Get("/api/me", middleware.RequireAuthAPI(authHandler.Me))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	staticDir := "web/static"
	if _, err := os.Stat(staticDir); err == nil {
		fileServer := http.FileServer(http.Dir(staticDir))
		r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	}

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s (%s)", cfg.App.Port, cfg.App.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}

func runMigrations(dbURL string) error {
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate close source: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate close db: %w", dbErr)
	}

	return nil
}

func loadTemplates(root string) *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"upper": strings.ToUpper,
		"seq": func(n int) []int {
			r := make([]int, n)
			for i := 0; i < n; i++ {
				r[i] = i
			}
			return r
		},
	})

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		name := filepath.ToSlash(rel)

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("warning: could not read template %s: %v", path, err)
			return nil
		}

		tmpl = template.Must(tmpl.New(name).Parse(string(content)))
		return nil
	})

	return tmpl
}
