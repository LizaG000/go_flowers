package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RabbitRequest struct {
	ID             uuid.UUID       `json:"id"`
	Version        string          `json:"version"`
	Action         string          `json:"action"`
	Data           json.RawMessage `json:"data"`
	Auth           string          `json:"auth"`
	IdempotencyKey uuid.UUID       `json:"idempotency_key"`
}

type RabbitResponse struct {
	CorrelationID uuid.UUID       `json:"correlation_id"`
	Status        string          `json:"status"`
	Data          json.RawMessage `json:"data"`
	Error         string          `json:"error,omitempty"`
}
