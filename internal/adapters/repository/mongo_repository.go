package repository

import (
	"context"
	"encoding/json"
	"strings"

	"mongodb-api/internal/core/domain"
	"mongodb-api/internal/core/ports"
	"mongodb-api/internal/infrastructure/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoRepository struct {
	client *db.MongoClient
}

func NewMongoRepository(client *db.MongoClient) ports.DocumentRepository {
	return &MongoRepository{client: client}
}

func (r *MongoRepository) idFilter(idStr string) bson.M {
	filters := []bson.M{
		{"_id": idStr},
		{"id": idStr},
	}
	if objID, err := bson.ObjectIDFromHex(idStr); err == nil {
		filters = append(filters, bson.M{"_id": objID})
	}
	return bson.M{"$or": filters}
}

func (r *MongoRepository) FindAll(ctx context.Context, database, collection string, filter domain.QueryFilter) ([]domain.Document, error) {
	col := r.client.GetCollection(database, collection)

	queryFilter := bson.M{}
	if filter.RawFilter != "" {
		var parsedFilter bson.M
		if err := json.Unmarshal([]byte(filter.RawFilter), &parsedFilter); err == nil {
			queryFilter = parsedFilter
		}
	}

	findOpts := options.Find()
	if filter.Limit > 0 {
		findOpts.SetLimit(filter.Limit)
	}
	if filter.Skip >= 0 {
		findOpts.SetSkip(filter.Skip)
	}
	if filter.SortField != "" {
		sortOrder := 1
		if filter.SortOrder == -1 {
			sortOrder = -1
		}
		findOpts.SetSort(bson.D{{Key: filter.SortField, Value: sortOrder}})
	}

	cursor, err := col.Find(ctx, queryFilter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.Document
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if results == nil {
		results = []domain.Document{}
	}

	return results, nil
}

func (r *MongoRepository) FindByID(ctx context.Context, database, collection, id string) (domain.Document, error) {
	col := r.client.GetCollection(database, collection)
	filter := r.idFilter(id)

	var doc domain.Document
	err := col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments || strings.Contains(err.Error(), "no documents") {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return doc, nil
}

func (r *MongoRepository) InsertOne(ctx context.Context, database, collection string, doc domain.Document) (any, error) {
	col := r.client.GetCollection(database, collection)
	if _, ok := doc["_id"]; !ok {
		if idVal, ok := doc["id"]; ok && idVal != "" {
			doc["_id"] = idVal
		}
	}
	res, err := col.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (r *MongoRepository) InsertMany(ctx context.Context, database, collection string, docs []domain.Document) ([]any, error) {
	col := r.client.GetCollection(database, collection)
	insertDocs := make([]any, len(docs))
	for i, d := range docs {
		if _, ok := d["_id"]; !ok {
			if idVal, ok := d["id"]; ok && idVal != "" {
				d["_id"] = idVal
			}
		}
		insertDocs[i] = d
	}

	res, err := col.InsertMany(ctx, insertDocs)
	if err != nil {
		return nil, err
	}

	return res.InsertedIDs, nil
}

func (r *MongoRepository) ReplaceOne(ctx context.Context, database, collection, id string, doc domain.Document) (int64, int64, error) {
	col := r.client.GetCollection(database, collection)
	filter := r.idFilter(id)

	delete(doc, "_id")

	res, err := col.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return 0, 0, err
	}

	return res.MatchedCount, res.ModifiedCount, nil
}

func (r *MongoRepository) UpdateOne(ctx context.Context, database, collection, id string, updates domain.Document) (int64, int64, error) {
	col := r.client.GetCollection(database, collection)
	filter := r.idFilter(id)

	delete(updates, "_id")

	updateQuery := bson.M{"$set": updates}
	res, err := col.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return 0, 0, err
	}

	return res.MatchedCount, res.ModifiedCount, nil
}

func (r *MongoRepository) DeleteOne(ctx context.Context, database, collection, id string) (int64, error) {
	col := r.client.GetCollection(database, collection)
	filter := r.idFilter(id)

	res, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (r *MongoRepository) Ping(ctx context.Context) error {
	if r.client == nil || r.client.Client == nil {
		return mongo.ErrClientDisconnected
	}
	return r.client.Client.Ping(ctx, nil)
}
