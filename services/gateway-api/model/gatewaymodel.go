package model

import (
	"context"
	"net/http"
)

// ===== Models =====
type RawRequestData struct {
	Ctx     context.Context
	Method  string
	Path    string
	Header  http.Header
	Body    []byte
	IP      string
	Topic   string
	Token   string
	ReplyCh chan GatewayResult
}
type GatewayResult struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type GatewayModel struct {
	RequestQueue chan RawRequestData
}

func NewGatewayModel() *GatewayModel {
	return &GatewayModel {
		RequestQueue : make(chan RawRequestData, 1024)
	}
}

