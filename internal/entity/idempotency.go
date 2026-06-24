package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Idempotency struct {
	Key          uuid.UUID       `json:"key" db:"key"`
	Status       string          `json:"status" db:"status"`
	ResponseCode int             `json:"response_code" db:"response_code"`
	ResponseBody json.RawMessage `json:"response_body" db:"response_body"`
	PayloadHash  string          `json:"payload_hash" db:"payload_hash"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateIdempotency struct {
	Key          uuid.UUID       `json:"key" db:"key"`
	Status       string          `json:"status" db:"status"`
	ResponseCode int             `json:"response_code" db:"response_code"`
	ResponseBody json.RawMessage `json:"response_body" db:"response_body"`
	PayloadHash  string          `json:"payload_hash" db:"payload_hash"`
}
