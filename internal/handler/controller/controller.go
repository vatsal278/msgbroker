package controller

import (
	"bytes"
	"github.com/vatsal278/msgbroker/internal/pkg/parser"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/handler"
	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
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
		err := parseRequest.ParseAndValidateRequest(r.Body, &publisher)
		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.IncompleteData, nil, &model.Response{})
			log.Println(err.Error())
			return
		}

		publisher.Id = uuid.New().String()
		x, ok := m.messageBroker.PubM[publisher.Channel]
		log.Print(x)
		if !ok {
			x = make(map[string]struct{})
			log.Print(publisher.Id)
		}
		x[publisher.Id] = struct{}{}
		m.messageBroker.PubM[publisher.Channel] = x
		log.Print(m.messageBroker.PubM)

		responseWriter.ResponseWriter(w, http.StatusCreated, constants.PublisherRegistration, map[string]interface{}{
			"id": publisher.Id,
		}, &model.Response{})
		log.Print(constants.PublisherRegistration)
	}
}

func (m *models) RegisterSubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscriber model.Subscriber
		err := parseRequest.ParseAndValidateRequest(r.Body, &subscriber)
		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.IncompleteData, nil, &model.Response{})
			log.Println(err.Error())
			return
		}
		log.Print(subscriber)
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
		responseWriter.ResponseWriter(w, http.StatusCreated, constants.SubscriberRegistration, nil, &model.Response{})
		log.Print(constants.SubscriberRegistration)
	}
}

func (m *models) PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		err := parseRequest.ParseAndValidateRequest(r.Body, &updates)

		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.IncompleteData, nil, &model.Response{})
			log.Println(err.Error())
			return
		}
		_, err = uuid.Parse(updates.Publisher.Id)

		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.InvalidUUID, nil, &model.Response{})
			log.Println(err.Error())
			return
		}

		pubm := m.messageBroker.PubM[updates.Publisher.Channel]
		_, ok := pubm[updates.Publisher.Id]
		if !ok {
			responseWriter.ResponseWriter(w, http.StatusNotFound, constants.PublisherNotFound, nil, &model.Response{})
			log.Println(constants.PublisherNotFound)
			return
		}

		for _, v := range m.messageBroker.SubM[updates.Publisher.Channel] {
			go func(v model.Subscriber) {
				client := http.Client{
					Timeout: time.Duration(2 * time.Second),
				}
				method := v.CallBack.HttpMethod
				url := v.CallBack.CallbackUrl
				reqBody := []byte(updates.Update)
				if v.CallBack.PublicKey != "" {

					PublicKey := v.CallBack.PublicKey
					PubKey, err := crypt.PEMStrAsKey(PublicKey)
					if err != nil {
						log.Print(err.Error())
						return
					}
					a, err := crypt.RsaOaepEncrypt(updates.Update, *PubKey)
					if err != nil {
						log.Print(err.Error())
						return
					}

					reqBody = []byte(a)
				}
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
		responseWriter.ResponseWriter(w, http.StatusOK, constants.NotifiedSub, nil, &model.Response{})
		log.Println(constants.NotifiedSub)
	}
}
func (m *models) NoRouteFound() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseWriter.ResponseWriter(w, http.StatusNotFound, constants.NoRoute, nil, &model.Response{})
		log.Print(constants.NoRoute)
	}
}
