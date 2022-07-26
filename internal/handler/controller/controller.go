package article_controller

import (
	"bytes"
	"encoding/json"
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
		err := parser.ParseResponse(r.Body, &publisher)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Print(err.Error())
			}
		}
		err = parser.ValidateRequest(&publisher)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusBadRequest, constants.Incomplete_Data, nil))
			if err != nil {
				log.Println(err.Error())
			}
		}
		x, ok := MessageBroker.PubM[publisher.Channel]
		if !ok {
			x = make(map[string]struct{})
			x[publisher.Name] = struct{}{}
		}
		MessageBroker.PubM[publisher.Channel] = x

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusCreated, "Successfully Registered as publisher to the channel", nil))
		if err != nil {
			log.Println(err.Error())
		}
		log.Print("Successfully Registered as publisher to the channel")
	}
}

func RegisterSubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscriber model.Subscriber
		w.Header().Set("Content-Type", "application/json")
		err := parser.ParseResponse(r.Body, &subscriber)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Print(err.Error())
			}
		}
		err = parser.ValidateRequest(&subscriber)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusBadRequest, constants.Incomplete_Data, nil))
			if err != nil {
				log.Println(err.Error())
			}
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
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusCreated, "Successfully Subscribed to the channel", nil))
		if err != nil {
			log.Println(err.Error())
		}
		log.Print("Successfully Subscribed to the channel")
	}
}

func PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		w.Header().Set("Content-Type", "application/json")
		err := parser.ParseResponse(r.Body, &updates)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Print(err.Error())
			}
		}
		err = parser.ValidateRequest(&updates)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusBadRequest, constants.Incomplete_Data, nil))
			if err != nil {
				log.Println(err.Error())
			}
		}

		pubm := MessageBroker.PubM[updates.Publisher.Channel]
		_, ok := pubm[updates.Publisher.Name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusNotFound, "No publisher found with the specified name for specified channel", nil))
			if err != nil {
				log.Println(err.Error())
			}
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
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(parser.Response_Writer(http.StatusOK, "notified all subscriber", nil))
		if err != nil {
			log.Println(err.Error())
		}
	}
}
