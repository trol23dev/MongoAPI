package ports

import (
	"context"

	"mongodb-api/internal/core/domain"
)

// DocumentRepository defines the driven (outbound) port for database operations.
type DocumentRepository interface {
	FindAll(ctx context.Context, database, collection string, filter domain.QueryFilter) ([]domain.Document, error)
	FindByID(ctx context.Context, database, collection, id string) (domain.Document, error)
	InsertOne(ctx context.Context, database, collection string, doc domain.Document) (any, error)
	InsertMany(ctx context.Context, database, collection string, docs []domain.Document) ([]any, error)
	ReplaceOne(ctx context.Context, database, collection, id string, doc domain.Document) (int64, int64, error)
	UpdateOne(ctx context.Context, database, collection, id string, updates domain.Document) (int64, int64, error)
	DeleteOne(ctx context.Context, database, collection, id string) (int64, error)
	Ping(ctx context.Context) error
}
