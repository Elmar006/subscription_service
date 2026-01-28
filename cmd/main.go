package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Elmar006/subscription_service/internal/config"
	"github.com/Elmar006/subscription_service/internal/db"
	"github.com/Elmar006/subscription_service/internal/handler"
	"github.com/Elmar006/subscription_service/internal/repository"
	"github.com/Elmar006/subscription_service/logger"
)

func main() {
	cfg := config.Load()
	logger.Init()
	log := logger.L()

	database := db.Connect(cfg)
	defer database.Close()

	repo := repository.NewSubscriptionRepo(database)
	handler := handler.NewSubscriptionHandler(repo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Post("/subscriptions", handler.CreateSubscription)
	r.Get("/subscriptions/{id}", handler.GetByIDSubscription)
	r.Put("/subscriptions/{id}", handler.UpdateByIDSubscription)
	r.Delete("/subscriptions/{id}", handler.DeleteSubscription)
	r.Get("/subscriptions", handler.GetSubscription)
	r.Get("/subscriptions/total", handler.GetSubscriptionTotal)

	log.Infof("Server started on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}
