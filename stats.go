package main

import (
	"time"
)

// CompletedStats is a statistics data structure for compseted backup
type CompletedStats struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Duration  int64     `json:"duration"`
}
