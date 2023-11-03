package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"seat-reservation/pkg/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize config
	var cfg Config
	if err := cfg.Load(); err != nil {
		log.Printf("failed to read config\n")
	}

	ctx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Printf("main panicked because: %v\n", err)
			} else {
				log.Printf("main panicked because: %v\n", r)
			}
		}
	}()
	err := run(ctx, cfg)
	if err != nil {
		log.Printf("server error: %v\n", err)
	}
}

func run(ctx context.Context, cfg Config) error {
	router := gin.Default()

	middleware, err := handler.NewMiddleware(ctx, &cfg.Config.DB)
	if err != nil {
		return err
	}
	router.Use(middleware.AuthMiddleware(cfg.Config.Auth))

	if err := handler.InitHandlers(ctx, middleware, router, cfg.Config); err != nil {
		return err
	}

	// run the server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}
	fmt.Printf("server is running on port %d\n", cfg.Server.Port)

	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
