package domain

import "time"

type Document = map[string]any

type InsertResult struct {
	ID any `json:"id,omitempty"`
}

type InsertManyResult struct {
	IDs   []any `json:"ids"`
	Count int   `json:"count"`
}

type UpdateResult struct {
	MatchedCount  int64 `json:"matchedCount"`
	ModifiedCount int64 `json:"modifiedCount"`
}

type DeleteResult struct {
	DeletedCount int64 `json:"deletedCount"`
}

type HealthStatus struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Time     string `json:"time"`
}

func NewHealthStatus(dbStatus string) HealthStatus {
	return HealthStatus{
		Status:   "ok",
		Database: dbStatus,
		Time:     time.Now().Format(time.RFC3339),
	}
}
