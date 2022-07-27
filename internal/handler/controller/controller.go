package controller

import (
	"bytes"
	"log"
	"net/http"
	"reflect"
	"time"

	controllerInterface "github.com/vatsal278/msgbroker/internal/handler"
	"github.com/vatsal278/msgbroker/internal/model"
	parser "github.com/vatsal278/msgbroker/internal/pkg/parser"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
)

type Models struct {
	publisher  model.Publisher
	subscriber model.Subscriber
	updates    model.Updates
}

var SubscriberMap = map[string][]model.Subscriber{}
var PublisherMap = map[string]map[string]struct{}{}

var MessageBroker = model.MessageBroker{
	SubM: SubscriberMap,
	PubM: PublisherMap,
}

func NewController() controllerInterface.IController {
	return &Models{
		publisher:  model.Publisher{},
		subscriber: model.Subscriber{},
		updates:    model.Updates{},
	}
}

func (m Models) RegisterPublisher() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parser.Parse(w, r.Body, m.publisher)
		x, ok := MessageBroker.PubM[m.publisher.Channel]
		if !ok {
			x = make(map[string]struct{})
			x[m.publisher.Name] = struct{}{}
		}
		MessageBroker.PubM[m.publisher.Channel] = x
		responseWriter.ResponseWriter(w, http.StatusCreated, "Successfully Registered as publisher to the channel", nil, &model.Response{})
		log.Print("Successfully Registered as publisher to the channel")
	}
}

func (m Models) RegisterSubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parser.Parse(w, r.Body, m.publisher)

		go func(s model.Subscriber) {
			MessageBroker.Lock()
			defer MessageBroker.Unlock()
			subs := MessageBroker.SubM[s.Channel]

			for _, v := range subs {
				if reflect.DeepEqual(v, s) {
					return
				}
			}
			subs = append(subs, s)
			log.Printf("subscriber added %+v", s)
			MessageBroker.SubM[s.Channel] = subs

		}(m.subscriber)
		responseWriter.ResponseWriter(w, http.StatusCreated, "Successfully Registered as publisher to the channel", nil, &model.Response{})
		log.Print("Successfully Subscribed to the channel")
	}
}

func (m Models) PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parser.Parse(w, r.Body, m.publisher)

		pubm := MessageBroker.PubM[m.updates.Publisher.Channel]
		_, ok := pubm[m.updates.Publisher.Name]
		if !ok {
			responseWriter.ResponseWriter(w, http.StatusNotFound, "No publisher found with the specified name for specified channel", nil, &model.Response{})
			log.Println("No publisher found with the specified name for specified channel")
			return
		}

		for _, v := range MessageBroker.SubM[m.updates.Publisher.Channel] {
			go func(v model.Subscriber) {
				reqBody := []byte(m.updates.Update)

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
