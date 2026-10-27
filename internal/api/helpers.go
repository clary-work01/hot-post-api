package api

import (
	"fmt"
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/domain"
	"github.com/google/uuid"
)

func parseID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {

		return uuid.Nil, fmt.Errorf("parse id %w", domain.ErrInvalidInput)
	}
	return id, nil
}
