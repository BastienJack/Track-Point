package db

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Event struct {
	gorm.Model
	EventName   string          `gorm:"type:varchar(256);not null;index:idx_event_name;" json:"event_name,omitempty"`
	EventParams json.RawMessage `gorm:"type:json;not null;" json:"event_params,omitempty"`
}

func (Event) TableName() string {
	return "Event"
}

func AddEvent(ctx context.Context, event *Event) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(db *gorm.DB) error {
		if err := db.Create(event).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func DeleteEventById(ctx context.Context, eventId uint) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(db *gorm.DB) error {
		// query event
		event := new(Event)
		if err := db.Where("id = ?", eventId).First(&event).Error; err != nil {
			return err
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		err := db.Delete(&event).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func GetAllEvents(ctx context.Context) ([]*Event, error) {
	var events []*Event
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Event{}).
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func GetPageEvents(ctx context.Context, offset int, limit int) ([]*Event, error) {
	var events []*Event
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Event{}).
		Offset(offset).Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}
