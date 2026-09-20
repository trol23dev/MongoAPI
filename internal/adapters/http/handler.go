package http

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"mongodb-api/internal/core/domain"
	"mongodb-api/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct {
	useCase ports.DocumentUseCase
}

func NewDocumentHandler(useCase ports.DocumentUseCase) *DocumentHandler {
	return &DocumentHandler{useCase: useCase}
}

func (h *DocumentHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Documento no encontrado",
		})
	case errors.Is(err, domain.ErrInvalidID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El ID proporcionado es inválido o está vacío",
		})
	case errors.Is(err, domain.ErrEmptyBody):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El cuerpo de la petición está vacío",
		})
	case errors.Is(err, domain.ErrDatabaseEmpty):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El nombre de la base de datos es requerido",
		})
	case errors.Is(err, domain.ErrCollectionEmpty):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El nombre de la colección es requerido",
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error interno del servidor",
			"details": err.Error(),
		})
	}
}

// GetAll handles GET /api/:database/:collection
func (h *DocumentHandler) GetAll(c *fiber.Ctx) error {
	database := c.Params("database")
	collection := c.Params("collection")

	filter := domain.QueryFilter{
		RawFilter: c.Query("filter"),
		SortField: c.Query("sort"),
		SortOrder: 1,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if skipStr := c.Query("skip"); skipStr != "" {
		if skip, err := strconv.ParseInt(skipStr, 10, 64); err == nil && skip >= 0 {
			filter.Skip = skip
		}
	}

	if strings.HasPrefix(filter.SortField, "-") {
		filter.SortOrder = -1
		filter.SortField = strings.TrimPrefix(filter.SortField, "-")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	docs, err := h.useCase.GetAll(ctx, database, collection, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(docs)
}

// GetByID handles GET /api/:database/:collection/:id
func (h *DocumentHandler) GetByID(c *fiber.Ctx) error {
	database := c.Params("database")
	collection := c.Params("collection")
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	doc, err := h.useCase.GetByID(ctx, database, collection, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(doc)
}

// Create handles POST /api/:database/:collection
func (h *DocumentHandler) Create(c *fiber.Ctx) error {
	database := c.Params("database")
	collection := c.Params("collection")

	bodyBytes := c.Body()
	trimmed := strings.TrimSpace(string(bodyBytes))

	if len(trimmed) == 0 {
		return h.handleError(c, domain.ErrEmptyBody)
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if strings.HasPrefix(trimmed, "[") {
		var docs []domain.Document
		if err := json.Unmarshal(bodyBytes, &docs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Formato de array JSON inválido",
				"details": err.Error(),
			})
		}

		result, err := h.useCase.CreateMany(ctx, database, collection, docs)
		if err != nil {
			return h.handleError(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":       "Datos insertados correctamente",
			"insertedIds":   result.IDs,
			"insertedCount": result.Count,
		})
	}

	var doc domain.Document
	if err := json.Unmarshal(bodyBytes, &doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Formato de objeto JSON inválido",
			"details": err.Error(),
		})
	}

	result, err := h.useCase.Create(ctx, database, collection, doc)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Data insertada correctamente",
		"insertedId": result.ID,
	})
}

// Replace handles PUT /api/:database/:collection/:id
func (h *DocumentHandler) Replace(c *fiber.Ctx) error {
	database := c.Params("database")
	collection := c.Params("collection")
	id := c.Params("id")

	var doc domain.Document
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cuerpo de la petición inválido",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	result, err := h.useCase.Replace(ctx, database, collection, id, doc)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Documento actualizado/reemplazado correctamente",
		"matchedCount":  result.MatchedCount,
		"modifiedCount": result.ModifiedCount,
	})
}

// UpdatePartial handles PATCH /api/:database/:collection/:id
func (h *DocumentHandler) UpdatePartial(c *fiber.Ctx) error {
	database := c.Params("database")
	collection := c.Params("collection")
	id := c.Params("id")

	var updates domain.Document
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cuerpo de la petición inválido",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	result, err := h.useCase.UpdatePartial(ctx, database, collection, id, updates)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Documento actualizado parcialmente con éxito",
		"matchedCount":  result.MatchedCount,
		"modifiedCount": result.ModifiedCount,
	})
}

// Delete handles DELETE /api/:database/:collection/:id o /api/:database/:collection?id=...
func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	database := c.Params("database")
	collection := c.Params("collection")
	id := c.Params("id")
	if id == "" {
		id = c.Query("id")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	result, err := h.useCase.Delete(ctx, database, collection, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Data eliminada correctamente",
		"deletedCount": result.DeletedCount,
	})
}

// HealthCheck handles GET /health
func (h *DocumentHandler) HealthCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	status := h.useCase.HealthCheck(ctx)
	return c.Status(fiber.StatusOK).JSON(status)
}
