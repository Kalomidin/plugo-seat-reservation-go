package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepository struct {
	*gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db,
	}
}

func (repo *eventRepository) CreateEvent(ctx context.Context, event *Event) error {
	res := repo.Create(event)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo *eventRepository) GetEventWithSeats(ctx context.Context, id uuid.UUID) (*Event, error) {
	var event Event
	resp := repo.DB.Preload("Seats").Where("id = ?", id).First(&event)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return &event, nil
}

func (repo *eventRepository) Migrate() error {
	return repo.AutoMigrate(&Event{})
}

func (repo *eventRepository) MigrateDown() error {
	return repo.Migrator().DropTable(&Event{})
}
