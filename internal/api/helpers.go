package api

import (
	"net/http"

	"github.com/google/uuid"
)

func parseID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil || id == uuid.Nil {
		return uuid.Nil, err
	}
	return id, nil
}
