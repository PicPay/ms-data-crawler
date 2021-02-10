package assembler

import (
	"time"
)

type AssembledScreen struct {
	Identifier string      `json:"identifier"`
	Version    int         `json:"version"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Body       interface{} `json:"body"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
