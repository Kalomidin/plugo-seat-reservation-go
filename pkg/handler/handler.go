package handler

import (
	"context"
	"fmt"
	"seat-reservation/common"
	"seat-reservation/pkg/config"
	auth_handler "seat-reservation/pkg/handler/auth"
	event_handler "seat-reservation/pkg/handler/event"
	user_handler "seat-reservation/pkg/handler/user"
	auth_manager "seat-reservation/pkg/manager/auth"
	auth_repo "seat-reservation/pkg/manager/auth/repository"
	event_manager "seat-reservation/pkg/manager/event"
	"seat-reservation/pkg/manager/middleware"
	user_manager "seat-reservation/pkg/manager/user"
	user_repo "seat-reservation/pkg/manager/user/repository"

	event_repo "seat-reservation/pkg/manager/event/repository"

	"seat-reservation/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewMiddleware(ctx context.Context, dbConfig postgres.ConfigPostgres) (common.Middleware, error) {
	db, err := postgres.InitDB(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	userRepo := user_repo.NewRepository(ctx, db)
	userPort := user_manager.NewUserPort(userRepo)
	return middleware.NewMiddleware(userPort), nil
}

func InitHandlers(
	ctx context.Context,
	middleware common.Middleware,
	ginEngine *gin.Engine,
	cfg config.Config) error {

	db, err := postgres.InitDB(ctx, &cfg.DB)
	if err != nil {
		return err
	}

	userRepo := user_repo.NewRepository(ctx, db)
	authRepo := auth_repo.NewRepository(ctx, db)

	eventRepo := event_repo.NewEventRepository(db)
	seatRepo := event_repo.NewSeatRepository(db)
	reservationRepo := event_repo.NewReservationRepository(db)

	// run migrations
	if err := runMigrations(db, &cfg.DB,
		userRepo,
		authRepo,
		eventRepo,
		seatRepo,
		reservationRepo,
	); err != nil {
		return err
	}

	userPort := user_manager.NewUserPort(userRepo)
	authManager := auth_manager.NewAuthManager(authRepo, userPort, cfg.Auth)
	authHandler := auth_handler.NewHandler(authManager)
	authHttpHandler := auth_handler.NewHttpHandler(ctx, authHandler, middleware)
	authHttpHandler.Init(ctx, ginEngine)

	userManager := user_manager.NewUserManager(userRepo)
	userHandler := user_handler.NewHandler(userManager)
	userHttpHandler := user_handler.NewHttpHandler(ctx, userHandler, middleware)
	userHttpHandler.Init(ctx, ginEngine)

	eventManager := event_manager.NewEventManager(eventRepo, seatRepo, reservationRepo)
	eventHandler := event_handler.NewHandler(eventManager)
	eventHttpHandler := event_handler.NewHttpHandler(ctx, eventHandler, middleware)
	eventHttpHandler.Init(ctx, ginEngine)

	return nil
}

func runMigrations(db *gorm.DB,
	dbConfig postgres.ConfigPostgres,
	userRepo user_repo.Repository,
	authRepo auth_repo.Repository,
	eventRepo event_repo.EventRepository,
	seatRepo event_repo.SeatRepository,
	reservationRepo event_repo.ReservationRepository,
) error {
	if resp := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", dbConfig.GetSchema())); resp.Error != nil {
		return resp.Error
	}

	if err := userRepo.Migrate(); err != nil {
		return err
	}
	if err := authRepo.Migrate(); err != nil {
		return err
	}

	if err := eventRepo.Migrate(); err != nil {
		return err
	}
	if err := seatRepo.Migrate(); err != nil {
		return err
	}
	if err := reservationRepo.Migrate(); err != nil {
		return err
	}

	return nil
}
