package ports

import (
	"context"

	"mongodb-api/internal/core/domain"
)

// DocumentUseCase defines the driving (inbound) port for document business operations.
type DocumentUseCase interface {
	GetAll(ctx context.Context, database, collection string, filter domain.QueryFilter) ([]domain.Document, error)
	GetByID(ctx context.Context, database, collection, id string) (domain.Document, error)
	Create(ctx context.Context, database, collection string, doc domain.Document) (*domain.InsertResult, error)
	CreateMany(ctx context.Context, database, collection string, docs []domain.Document) (*domain.InsertManyResult, error)
	Replace(ctx context.Context, database, collection, id string, doc domain.Document) (*domain.UpdateResult, error)
	UpdatePartial(ctx context.Context, database, collection, id string, updates domain.Document) (*domain.UpdateResult, error)
	Delete(ctx context.Context, database, collection, id string) (*domain.DeleteResult, error)
	HealthCheck(ctx context.Context) domain.HealthStatus
}
