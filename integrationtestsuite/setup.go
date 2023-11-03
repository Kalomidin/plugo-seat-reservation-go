package integrationtestsuite

import (
	"context"
	"seat-reservation/common"
	"seat-reservation/jwt"
	"seat-reservation/pkg/config"
	"seat-reservation/pkg/handler"
	"seat-reservation/postgres"
	"time"

	"gorm.io/gorm"

	auth_manager "seat-reservation/pkg/manager/auth"
	auth_repo "seat-reservation/pkg/manager/auth/repository"
	event_manager "seat-reservation/pkg/manager/event"
	"seat-reservation/pkg/manager/port"
	user_manager "seat-reservation/pkg/manager/user"
	user_repo "seat-reservation/pkg/manager/user/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	event_repo "seat-reservation/pkg/manager/event/repository"
)

type DBIntegrationSuite struct {
	suite.Suite
	Context         context.Context
	DB              *gorm.DB
	UserRepo        user_repo.Repository
	AuthRepo        auth_repo.Repository
	EventRepo       event_repo.EventRepository
	SeatRepo        event_repo.SeatRepository
	ReservationRepo event_repo.ReservationRepository
	UserPort        port.UserPort
	User            user_manager.User
	Event           event_manager.Event
}

func (r *DBIntegrationSuite) SetupSuite() {
	r.Context = context.Background()
	dbConfig := getTestDBConfig()
	db, err := postgres.InitDB(r.Context, &dbConfig)
	require.NoError(r.T(), err)

	r.UserRepo = user_repo.NewRepository(r.Context, db)
	r.AuthRepo = auth_repo.NewRepository(r.Context, db)
	r.EventRepo = event_repo.NewEventRepository(db)
	r.SeatRepo = event_repo.NewSeatRepository(db)
	r.ReservationRepo = event_repo.NewReservationRepository(db)
	r.UserPort = user_manager.NewUserPort(r.UserRepo)

	r.DB = db

	r.SetupTest()

	r.User = r.NewUser()
	r.Event = r.NewEvent(nil)
}

func (r *DBIntegrationSuite) SetupTest() {
	dbConfig := getTestDBConfig()
	require.NoError(r.T(), handler.RunMigrations(
		r.DB,
		&dbConfig,
		r.UserRepo,
		r.AuthRepo,
		r.EventRepo,
		r.SeatRepo,
		r.ReservationRepo,
	))
}

// func (r *DBIntegrationSuite) TearDownTest() {
// 	require.NoError(r.T(), func() error {
// 		if err := r.ReservationRepo.MigrateDown(); err != nil {
// 			return err
// 		}
// 		if err := r.SeatRepo.MigrateDown(); err != nil {
// 			return err
// 		}

// 		if err := r.EventRepo.MigrateDown(); err != nil {
// 			return err
// 		}
// 		if err := r.AuthRepo.MigrateDown(); err != nil {
// 			return err
// 		}
// 		if err := r.UserRepo.MigrateDown(); err != nil {
// 			return err
// 		}
// 		return nil
// 	}())

// }

// func (r *DBIntegrationSuite) TearDownSuite() {
// 	sqlDB, err := r.DB.DB()
// 	r.TearDownTest()

// 	require.NoError(r.T(), err)
// 	require.NoError(r.T(), sqlDB.Close())
// }

func (r *DBIntegrationSuite) EventManager() event_manager.EventManager {
	return event_manager.NewEventManager(
		r.EventRepo,
		r.SeatRepo,
		r.ReservationRepo,
		event_repo.NewTransactionDB(r.DB),
	)
}

func (r *DBIntegrationSuite) UserManager() user_manager.UserManager {
	return user_manager.NewUserManager(
		r.UserRepo,
	)
}

func (r *DBIntegrationSuite) AuthManager() auth_manager.AuthManager {
	return auth_manager.NewAuthManager(
		r.AuthRepo,
		r.UserPort,
		getTestAuthConfig(),
	)
}

func (r *DBIntegrationSuite) NewUser() user_manager.User {
	authManager := r.AuthManager()
	req := auth_manager.SignupOrLoginRequest{
		UserName: common.RandomString(20),
		Password: common.RandomString(20),
	}
	resp, err := authManager.SignupOrLogin(r.Context, req)
	require.NoError(r.T(), err)

	userManager := r.UserManager()
	user, err := userManager.GetUser(r.Context, resp.Id)
	require.NoError(r.T(), err)

	return *user
}

func (r *DBIntegrationSuite) NewEvent(userId *uuid.UUID) event_manager.Event {
	creatorID := r.User.ID
	if userId != nil {
		creatorID = *userId
	}

	eventManager := r.EventManager()
	req := event_manager.CreateEventRequest{
		Name:      common.RandomString(20),
		SeatCount: 10,
		CreatorID: creatorID,
	}

	resp, err := eventManager.CreateEvent(r.Context, req)
	require.NoError(r.T(), err)

	return *resp
}

func getTestAuthConfig() jwt.Config {
	return jwt.Config{
		Issuer:                     "seat-reservation",
		Key:                        "key",
		KeyID:                      "key-id",
		ValidDuration:              1 * time.Hour,
		RefreshTokenExpiryDuration: 24 * time.Hour,
	}
}

func getTestDBConfig() config.Postgres {
	return config.Postgres{
		Address:              ":5432",
		Host:                 "localhost",
		Database:             "seat-reservation",
		Port:                 "5432",
		Username:             "root",
		Password:             "password",
		MaxConns:             10,
		MaxWaitForConnection: 5 * time.Second,
		SSLMode:              "disable",
		Schema:               "seat_reservation_test",
	}
}
