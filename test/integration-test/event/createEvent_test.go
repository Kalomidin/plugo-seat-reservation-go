package event_test

import (
	"context"
	"seat-reservation/integrationtestsuite"
	event_manager "seat-reservation/pkg/manager/event"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type createEventSuite struct {
	integrationtestsuite.DBIntegrationSuite
}

func TestCreateEventSuite(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	suite.Run(t, new(createEventSuite))
}

func (s *createEventSuite) TestCreateReservation() {

	eventManager := s.EventManager()
	ctx := context.Background()
	req := event_manager.CreateEventRequest{
		Name:      "Test Event",
		CreatorID: s.User.ID,
		SeatCount: 10,
	}

	resp, err := eventManager.CreateEvent(ctx, req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), req.Name, resp.Name)
	require.Equal(s.T(), req.CreatorID, resp.CreatorID)
	require.Equal(s.T(), req.SeatCount, len(resp.Seats))
}
