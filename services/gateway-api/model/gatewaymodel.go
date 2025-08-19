package model

import (
	"context"
	apis "gatewayapi/internal/api"
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
	RequestQueue  chan RawRequestData
	TopicAuthMap  map[string]bool
	ServiceGroups []apis.ServiceGroup
}

func NewGatewayModel() *GatewayModel {
	var gateway = &GatewayModel{}
	gateway.RequestQueue = make(chan RawRequestData, 1024)
	gateway.ServiceGroups = []apis.ServiceGroup{
		apis.AuthService,
		apis.UserService,
		apis.PostsService,
		apis.ReactionsService,
	}
	gateway.TopicAuthMap = make(map[string]bool)
	gateway.initTopicAuthMap()
	return gateway
}

func (g *GatewayModel) initTopicAuthMap() {
	for _, sg := range g.ServiceGroups {
		for _, ep := range sg.Endpoints {
			topic := sg.Name + "/" + ep.Name
			g.TopicAuthMap[topic] = ep.RequireAuth
		}
	}
}
