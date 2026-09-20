package services

import (
	"context"
	"strings"

	"mongodb-api/internal/core/domain"
	"mongodb-api/internal/core/ports"
)

type DocumentService struct {
	repo ports.DocumentRepository
}

func NewDocumentService(repo ports.DocumentRepository) ports.DocumentUseCase {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) validateTarget(database, collection string) error {
	if strings.TrimSpace(database) == "" {
		return domain.ErrDatabaseEmpty
	}
	if strings.TrimSpace(collection) == "" {
		return domain.ErrCollectionEmpty
	}
	return nil
}

func (s *DocumentService) GetAll(ctx context.Context, database, collection string, filter domain.QueryFilter) ([]domain.Document, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	return s.repo.FindAll(ctx, database, collection, filter)
}

func (s *DocumentService) GetByID(ctx context.Context, database, collection, id string) (domain.Document, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidID
	}
	return s.repo.FindByID(ctx, database, collection, id)
}

func (s *DocumentService) Create(ctx context.Context, database, collection string, doc domain.Document) (*domain.InsertResult, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	if len(doc) == 0 {
		return nil, domain.ErrEmptyBody
	}

	id, err := s.repo.InsertOne(ctx, database, collection, doc)
	if err != nil {
		return nil, err
	}

	return &domain.InsertResult{ID: id}, nil
}

func (s *DocumentService) CreateMany(ctx context.Context, database, collection string, docs []domain.Document) (*domain.InsertManyResult, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, domain.ErrEmptyBody
	}

	ids, err := s.repo.InsertMany(ctx, database, collection, docs)
	if err != nil {
		return nil, err
	}

	return &domain.InsertManyResult{
		IDs:   ids,
		Count: len(ids),
	}, nil
}

func (s *DocumentService) Replace(ctx context.Context, database, collection, id string, doc domain.Document) (*domain.UpdateResult, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidID
	}
	if len(doc) == 0 {
		return nil, domain.ErrEmptyBody
	}

	matched, modified, err := s.repo.ReplaceOne(ctx, database, collection, id, doc)
	if err != nil {
		return nil, err
	}
	if matched == 0 {
		return nil, domain.ErrNotFound
	}

	return &domain.UpdateResult{
		MatchedCount:  matched,
		ModifiedCount: modified,
	}, nil
}

func (s *DocumentService) UpdatePartial(ctx context.Context, database, collection, id string, updates domain.Document) (*domain.UpdateResult, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidID
	}
	if len(updates) == 0 {
		return nil, domain.ErrEmptyBody
	}

	matched, modified, err := s.repo.UpdateOne(ctx, database, collection, id, updates)
	if err != nil {
		return nil, err
	}
	if matched == 0 {
		return nil, domain.ErrNotFound
	}

	return &domain.UpdateResult{
		MatchedCount:  matched,
		ModifiedCount: modified,
	}, nil
}

func (s *DocumentService) Delete(ctx context.Context, database, collection, id string) (*domain.DeleteResult, error) {
	if err := s.validateTarget(database, collection); err != nil {
		return nil, err
	}
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidID
	}

	deleted, err := s.repo.DeleteOne(ctx, database, collection, id)
	if err != nil {
		return nil, err
	}
	if deleted == 0 {
		return nil, domain.ErrNotFound
	}

	return &domain.DeleteResult{
		DeletedCount: deleted,
	}, nil
}

func (s *DocumentService) HealthCheck(ctx context.Context) domain.HealthStatus {
	dbStatus := "connected"
	if s.repo == nil {
		dbStatus = "not initialized"
	} else if err := s.repo.Ping(ctx); err != nil {
		dbStatus = "disconnected: " + err.Error()
	}
	return domain.NewHealthStatus(dbStatus)
}
