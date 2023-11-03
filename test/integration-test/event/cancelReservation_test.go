package event_test

import (
	"seat-reservation/integrationtestsuite"
	event_manager "seat-reservation/pkg/manager/event"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
)

type cancelReservationSuite struct {
	integrationtestsuite.DBIntegrationSuite
}

func TestCancelReservationSuite(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	suite.Run(t, new(cancelReservationSuite))
}

func (s *cancelReservationSuite) TestCancelReservation() {
	event := s.Event
	user := s.User
	eventManager := s.EventManager()
	ctx := s.Context
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
		SeatID:  event.Seats[0].ID,
	}

	_, err := eventManager.CreateReservation(ctx, req)
	s.Require().NoError(err)

	cancelReq := event_manager.CancelReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
	}

	_, err = eventManager.CancelReservation(ctx, cancelReq)
	s.Require().NoError(err)
}

func (s *cancelReservationSuite) TestCancelReservation_InvalidUser() {
	event := s.Event
	user := s.User
	eventManager := s.EventManager()
	ctx := s.Context
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
		SeatID:  event.Seats[0].ID,
	}

	_, err := eventManager.CreateReservation(ctx, req)
	s.Require().NoError(err)

	cancelReq := event_manager.CancelReservationRequest{
		EventID: event.ID,
		UserID:  s.NewUser().ID,
	}

	_, err = eventManager.CancelReservation(ctx, cancelReq)
	s.Require().Error(err)
}

func (s *cancelReservationSuite) TestCancelReservation_InvalidEvent() {
	event := s.Event
	user := s.User
	eventManager := s.EventManager()
	ctx := s.Context
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
		SeatID:  event.Seats[0].ID,
	}

	_, err := eventManager.CreateReservation(ctx, req)
	s.Require().NoError(err)

	cancelReq := event_manager.CancelReservationRequest{
		EventID: s.NewEvent(nil).ID,
		UserID:  user.ID,
	}

	_, err = eventManager.CancelReservation(ctx, cancelReq)
	s.Require().Error(err)
}

/*
**Scenario 3**: User keeps on trying to cancel the item. Only one cancel should succeed and all others should fail.
 */
func (s *cancelReservationSuite) TestCancelReservation_MultipleCancel() {
	event := s.Event
	user := s.User
	eventManager := s.EventManager()
	ctx := s.Context
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
		SeatID:  event.Seats[0].ID,
	}

	_, err := eventManager.CreateReservation(ctx, req)
	s.Require().NoError(err)

	cancelReq := event_manager.CancelReservationRequest{
		EventID: event.ID,
		UserID:  user.ID,
	}

	workCount := 20
	work := sync.WaitGroup{}
	var successCount int32 = 0

	for i := 0; i < workCount; i++ {
		work.Add(1)
		go func() {
			if _, err = eventManager.CancelReservation(ctx, cancelReq); err == nil {
				atomic.AddInt32(&successCount, 1)
			}

			work.Done()
		}()
	}
	work.Wait()

	s.Require().Equal(int32(1), successCount)
}

/*
- **Scenario 4**: First user keeps on trying to cancel a reservation and at the same time second user
keeps on trying to reserve the same seat. Second user should succeed reserving after first user
reservation successfully canceled and final state of the reservation should be reserved not `canceled`.
*/
func (s *cancelReservationSuite) TestCancelReservation_MultipleCancelAndReserve() {
	event := s.Event
	firstUser := s.User
	eventManager := s.EventManager()
	ctx := s.Context
	req := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  firstUser.ID,
		SeatID:  event.Seats[0].ID,
	}

	_, err := eventManager.CreateReservation(ctx, req)
	s.Require().NoError(err)

	secondUser := s.NewUser()
	secondUserResReq := event_manager.CreateReservationRequest{
		EventID: event.ID,
		UserID:  secondUser.ID,
		SeatID:  event.Seats[0].ID,
	}

	cancelReq := event_manager.CancelReservationRequest{
		EventID: event.ID,
		UserID:  firstUser.ID,
	}

	workCount := 30
	work := sync.WaitGroup{}
	var cancelSuccessCount int32 = 0
	var createSuccessCount int32 = 0

	for i := 0; i < workCount; i++ {
		work.Add(2)
		go func() {
			if _, err = eventManager.CancelReservation(ctx, cancelReq); err == nil {
				atomic.AddInt32(&cancelSuccessCount, 1)
			}

			work.Done()
		}()

		go func() {
			if _, err = eventManager.CreateReservation(ctx, secondUserResReq); err == nil {
				atomic.AddInt32(&createSuccessCount, 1)
			}
			work.Done()
		}()

	}
	work.Wait()

	s.Require().Equal(int32(1), cancelSuccessCount)
	s.Require().Greater(createSuccessCount, int32(0))
}
