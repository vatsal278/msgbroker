package controllerInterface

import "net/http"

//go:generate mockgen --destination=./../../mocks/internal/mocks/mock_controller.go --package=mocks github.com/vatsal278/msgbroker/internal/handler IController

type IController interface {
	RegisterPublisher() func(w http.ResponseWriter, r *http.Request)
	RegisterSubscriber() func(w http.ResponseWriter, r *http.Request)
	PublishMessage() func(w http.ResponseWriter, r *http.Request)
	NoRouteFound() func(w http.ResponseWriter, r *http.Request)
}
