package models

import "time"

type AuditLog struct {
	ID         int
	UserID     int
	Action     string
	ObjectName string
	Timestamp  time.Time
	Details    string
	OldState   string
	NewState   string
}
