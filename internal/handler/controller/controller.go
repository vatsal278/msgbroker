package article_controller

import (
	"bytes"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/model"
	parser "github.com/vatsal278/msgbroker/internal/pkg/parser"
)

var SubscriberMap = map[string][]model.Subscriber{}
var PublisherMap = map[string]map[string]struct{}{}

var MessageBroker = model.MessageBroker{
	SubM: SubscriberMap,
	PubM: PublisherMap,
}

func RegisterPublisher() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var publisher model.Publisher
		w.Header().Set("Content-Type", "application/json")
		err := parser.ParseRequest(r.Body, &publisher)
		if err != nil {
			parser.Response_Writer(w, http.StatusInternalServerError, constants.Parse_Err, nil, model.Response{})
			log.Println(err.Error())
		}
		err = parser.ValidateRequest(&publisher)
		if err != nil {
			parser.Response_Writer(w, http.StatusBadRequest, constants.Incomplete_Data, nil, model.Response{})
			log.Println(err.Error())
		}
		x, ok := MessageBroker.PubM[publisher.Channel]
		if !ok {
			x = make(map[string]struct{})
			x[publisher.Name] = struct{}{}
		}
		MessageBroker.PubM[publisher.Channel] = x
		parser.Response_Writer(w, http.StatusCreated, "Successfully Registered as publisher to the channel", nil, model.Response{})
		log.Print("Successfully Registered as publisher to the channel")
	}
}

func RegisterSubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscriber model.Subscriber
		w.Header().Set("Content-Type", "application/json")
		err := parser.ParseRequest(r.Body, &subscriber)
		if err != nil {
			parser.Response_Writer(w, http.StatusInternalServerError, constants.Parse_Err, nil, model.Response{})
			log.Println(err.Error())
		}
		err = parser.ValidateRequest(&subscriber)
		if err != nil {
			parser.Response_Writer(w, http.StatusBadRequest, constants.Incomplete_Data, nil, model.Response{})
			log.Println(err.Error())
		}

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

		}(subscriber)
		parser.Response_Writer(w, http.StatusCreated, "Successfully Registered as publisher to the channel", nil, model.Response{})
		log.Print("Successfully Subscribed to the channel")
	}
}

func PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		w.Header().Set("Content-Type", "application/json")
		err := parser.ParseRequest(r.Body, &updates)
		if err != nil {
			parser.Response_Writer(w, http.StatusInternalServerError, constants.Parse_Err, nil, model.Response{})
			log.Println(err.Error())
		}
		err = parser.ValidateRequest(&updates)
		if err != nil {
			parser.Response_Writer(w, http.StatusBadRequest, constants.Incomplete_Data, nil, model.Response{})
			log.Println(err.Error())
		}

		pubm := MessageBroker.PubM[updates.Publisher.Channel]
		_, ok := pubm[updates.Publisher.Name]
		if !ok {
			parser.Response_Writer(w, http.StatusNotFound, "No publisher found with the specified name for specified channel", nil, model.Response{})
			log.Println("No publisher found with the specified name for specified channel")
			return
		}

		for _, v := range MessageBroker.SubM[updates.Publisher.Channel] {
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
		parser.Response_Writer(w, http.StatusOK, "notified all subscriber", nil, model.Response{})
		log.Println("notified all subscriber")
	}
}
