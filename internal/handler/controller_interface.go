package controller

import "net/http"

//all the methods associated with the endpoints are listed here

//go:generate mockgen --destination=./../../mocks/mock_controller.go --package=mocks challenge/internal/handler IController
type IController interface {
	RegisterPublisher() func(w http.ResponseWriter, r *http.Request)
	RegisterSubscriber() func(w http.ResponseWriter, r *http.Request)
	PublishMessage() func(w http.ResponseWriter, r *http.Request)
}
