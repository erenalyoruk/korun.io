package events

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	AccountRegisteredEvent EventType = "account.registered"
	AccountLoggedInEvent   EventType = "account.logged_in"
)

type Event struct {
	ID        string         `json:"id"`
	Type      EventType      `json:"type"`
	Source    string         `json:"source"`
	Data      map[string]any `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
	Version   string         `json:"version"`
}

func NewEvent(eventType EventType, source string, data map[string]any) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    source,
		Data:      data,
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

type AccountRegisteredData struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AccountLoggedInData struct {
	AccountID string    `json:"account_id"`
	Email     string    `json:"email"`
	LoginTime time.Time `json:"login_time"`
}
