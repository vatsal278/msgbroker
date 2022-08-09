package controller

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
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
			x[publisher.Id] = struct{}{}
		}
		x[publisher.Id] = struct{}{}
		m.messageBroker.PubM[publisher.Channel] = x
		log.Print(m.messageBroker.PubM)

		responseWriter.ResponseWriter(w, http.StatusCreated, "Successfully Registered as publisher to the channel", map[string]interface{}{
			"id": publisher.Id,
		}, &model.Response{})
		log.Print("Successfully Registered as publisher to the channel")
	}
}

func (m *models) RegisterSubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscriber model.Subscriber
		err := parser.ParseAndValidateRequest(r.Body, &subscriber)
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
		responseWriter.ResponseWriter(w, http.StatusCreated, "Successfully Registered as Subscriber to the channel", nil, &model.Response{})
		log.Print("Successfully Subscribed to the channel")
	}
}

func (m *models) PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		err := parser.ParseAndValidateRequest(r.Body, &updates)

		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.IncompleteData, nil, &model.Response{})
			log.Println(err.Error())
			return
		}
		_, err = uuid.Parse(updates.Publisher.Id)

		if err != nil {
			responseWriter.ResponseWriter(w, http.StatusBadRequest, "Invalid UUID", nil, &model.Response{})
			log.Println(err.Error())
			return
		}

		pubm := m.messageBroker.PubM[updates.Publisher.Channel]
		_, ok := pubm[updates.Publisher.Id]
		if !ok {
			responseWriter.ResponseWriter(w, http.StatusNotFound, "No publisher found with the specified name for specified channel", nil, &model.Response{})
			log.Println("No publisher found with the specified name for specified channel")
			return
		}
		PrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		PublicKey := PrivateKey.PublicKey
		//Encrypt Miryan Message
		message := updates.Update
		label := []byte("")
		hash := sha256.New()

		ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, &PublicKey, []byte(message), label)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		log.Print(PrivateKey)
		plainText, err := rsa.DecryptOAEP(hash, rand.Reader, PrivateKey, ciphertext, label)
		if err != nil {
			fmt.Println(err)
			return
		}
		updates.Update = string(ciphertext)

		log.Print(plainText)

		for _, v := range m.messageBroker.SubM[updates.Publisher.Channel] {
			go func(v model.Subscriber) {
				reqBody := []byte(updates.Update)
				//x := []byte(updates.Update.PublicKey)
				//reqBody = append(reqBody)
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
