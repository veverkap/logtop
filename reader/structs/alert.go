package structs

import (
	"time"
)

// Alert represents an alert
type Alert struct {
	Event     LogEvent
	HitCount  int
	Triggered time.Time
	Cleared   time.Time
}

// LoadAverageHitRate iterates through events
func LoadAverageHitRate(events []LogEvent, alertThresholdDuration int64) int64 {
	eventCount := int64(len(events))
	return eventCount / alertThresholdDuration
}
