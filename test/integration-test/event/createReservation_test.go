package event_test

import (
	"context"
	"seat-reservation/integrationtestsuite"
	event_manager "seat-reservation/pkg/manager/event"
	"sync"
	"sync/atomic"
	"testing"

	event_repo "seat-reservation/pkg/manager/event/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type createReservationSuite struct {
	integrationtestsuite.DBIntegrationSuite
}

func TestCreateReservationSuite(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	suite.Run(t, new(createReservationSuite))
}

func (s *createReservationSuite) TestCreateReservation() {
	event := s.Event
	user := s.User
	eventManager := s.EventManager()
	ctx := context.Background()
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
		SeatID:  event.Seats[0].ID,
	}

	resp, err := eventManager.CreateReservation(ctx, req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), req.EventID, resp.EventID)
	require.Equal(s.T(), req.SeatID, resp.Seat.ID)
}

/*
**Scenario 1**: One user send same seat reservation and they are proccessed concurrently. In this scenario, expected behavior is to return success
 */
func (s *createReservationSuite) TestCreateReservation_SameUserConcurrently() {
	event := s.Event
	user := s.User
	eventManager := s.EventManager()
	ctx := context.Background()
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
		SeatID:  event.Seats[0].ID,
	}
	var reservationId *uuid.UUID
	workCount := 10
	work := sync.WaitGroup{}

	for i := 0; i < workCount; i++ {
		work.Add(1)

		go func() {
			resp, err := eventManager.CreateReservation(ctx, req)
			require.NoError(s.T(), err)
			if reservationId == nil {
				reservationId = &resp.ReservationID
			} else {
				require.Equal(s.T(), *reservationId, resp.ReservationID)
			}
			work.Done()
			require.Equal(s.T(), req.EventID, resp.EventID)
			require.Equal(s.T(), req.SeatID, resp.Seat.ID)
		}()
	}
	work.Wait()
}

/*
**Scenario 2**: Two or more users wants to reserve same seat and they are executed concurrently. In this scenario, first executed will succeed and second will fail
 */
func (s *createReservationSuite) TestCreateReservation_DifferentUserConcurrently() {
	event := s.NewEvent(nil)
	eventManager := s.EventManager()
	seat := event.Seats[0]
	workCount := 20
	var successCount int32 = 0

	work := sync.WaitGroup{}
	for i := 0; i < workCount; i++ {
		work.Add(1)
		go func() {
			ctx := context.Background()
			user := s.NewUser()
			req := event_manager.CreateReservationRequest{
				EventID: event.ID,
				UserID:  user.ID,
				SeatID:  seat.ID,
			}

			if _, err := eventManager.CreateReservation(ctx, req); err == nil {
				atomic.AddInt32(&successCount, 1)
			} else {

			}
			work.Done()
		}()
	}
	work.Wait()

	require.Equal(s.T(), int32(1), successCount)

	// check the reservation for the event
	repoEvent, err := s.EventRepo.GetEventWithSeats(s.Context, event.ID)
	require.NoError(s.T(), err)
	for i, repoSeat := range repoEvent.Seats {
		if repoSeat.ID == seat.ID {
			require.Equal(s.T(), event_repo.SeatStatusReserved, repoSeat.Status, "seat %d should be reserved", i)
		} else {
			require.Equal(s.T(), event_repo.SeatStatusAvailable, repoSeat.Status, "seat %d should be available", i)
		}
	}
}
