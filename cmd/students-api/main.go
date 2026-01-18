package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cmanish049/students-api/internal/config"
	"github.com/cmanish049/students-api/internal/http/handlers/student"
	"github.com/cmanish049/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// setup database
	db, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	slog.Info("storage initialialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	defer db.Db.Close()
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	router.HandleFunc("POST /api/students", student.New(db))

	router.HandleFunc("GET /api/students/{id}", student.GetById(db))
	router.HandleFunc("GET /api/students", student.GetStudentList(db))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateStudent(db))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteStudent(db))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	// Graceful shutdown

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shoutdown successfully")
}
