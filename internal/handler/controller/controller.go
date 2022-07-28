package controller

import (
	"bytes"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/vatsal278/msgbroker/internal/constants"
	controllerInterface "github.com/vatsal278/msgbroker/internal/handler"
	"github.com/vatsal278/msgbroker/internal/model"
	parser "github.com/vatsal278/msgbroker/internal/pkg/parser"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
)

type models struct {
	messageBroker model.MessageBroker
}

func NewController() controllerInterface.IController {
	return &models{
		messageBroker: model.MessageBroker{
			SubM: map[string][]model.Subscriber{},
			PubM: map[string]map[string]struct{}{},
		},
	}
}

func (m *models) RegisterPublisher() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var publisher model.Publisher
		err := parser.ParseAndValidateRequest(r.Body, &publisher)
		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.Incomplete_Data, nil, &model.Response{})
			log.Println(err.Error())
			return
		}

		x, ok := m.messageBroker.PubM[publisher.Channel]
		if !ok {
			x = make(map[string]struct{})
			x[publisher.Name] = struct{}{}
		}
		m.messageBroker.PubM[publisher.Channel] = x
		responseWriter.ResponseWriter(w, http.StatusCreated, "Successfully Registered as publisher to the channel", nil, &model.Response{})
		log.Print("Successfully Registered as publisher to the channel")
	}
}

func (m *models) RegisterSubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscriber model.Subscriber
		err := parser.ParseAndValidateRequest(r.Body, &subscriber)
		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.Incomplete_Data, nil, &model.Response{})
			log.Println(err.Error())
			return
		}

		go func(s model.Subscriber) {
			m.messageBroker.Lock()
			defer m.messageBroker.Unlock()
			subs := m.messageBroker.SubM[s.Channel]

			for _, v := range subs {
				if reflect.DeepEqual(v, s) {
					return
				}
			}
			subs = append(subs, s)
			log.Printf("subscriber added %+v", s)
			m.messageBroker.SubM[s.Channel] = subs

		}(subscriber)
		responseWriter.ResponseWriter(w, http.StatusCreated, "Successfully Registered as publisher to the channel", nil, &model.Response{})
		log.Print("Successfully Subscribed to the channel")
	}
}

func (m *models) PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		err := parser.ParseAndValidateRequest(r.Body, &updates)
		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.Incomplete_Data, nil, &model.Response{})
			log.Println(err.Error())
			return
		}

		pubm := m.messageBroker.PubM[updates.Publisher.Channel]
		_, ok := pubm[updates.Publisher.Name]
		if !ok {
			responseWriter.ResponseWriter(w, http.StatusNotFound, "No publisher found with the specified name for specified channel", nil, &model.Response{})
			log.Println("No publisher found with the specified name for specified channel")
			return
		}

		for _, v := range m.messageBroker.SubM[updates.Publisher.Channel] {
			go func(v model.Subscriber) {
				reqBody := []byte(updates.Update)

				timeout := time.Duration(2 * time.Second)
				client := http.Client{
					Timeout: timeout,
				}
				method := v.CallBack.HttpMethod
				url := v.CallBack.CallbackUrl
				request, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
				if err != nil {
					log.Println(err.Error())
					return
				}
				request.Header.Set("Content-Type", "application/json")
				log.Printf("%+v \n", *request)
				client.Do(request)
			}(v)

		}
		responseWriter.ResponseWriter(w, http.StatusOK, "notified all subscriber", nil, &model.Response{})
		log.Println("notified all subscriber")
	}
}
func (m *models) NoRouteFound() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		responseWriter.ResponseWriter(w, http.StatusNotFound, "no route found", nil, &model.Response{})
		log.Print("No Route Found")
	}
}
