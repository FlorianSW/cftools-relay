package domain

import (
	"errors"
	"time"
)

var ErrCFToolsIdMissing = errors.New("CFTools ID is missing")

type EventHistory interface {
	Save(e Event) error
	FindWithin(eventType, cftoolsId string, within time.Duration) ([]Event, error)
}
