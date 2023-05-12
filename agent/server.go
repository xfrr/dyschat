package agent

import "context"

//go:generate moq -out mock/server.go -pkg mock . Server:ServerMock
type Server interface {
	Start(context.Context) error
}
