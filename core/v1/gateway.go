package v1

import "context"

//gateway.go
//go:generate go run github.com/golang/mock/mockgen -destination=../mocks/mock_http_client.go -package=mocks . HttpClient
type HttpClient interface {
	Send(ctx context.Context, configuration HttpConfiguration, request interface{}) (HttpResponse, error)
}
