package controller

import (
	"bytes"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/vatsal278/msgbroker/internal/constants"
	controllerInterface "github.com/vatsal278/msgbroker/internal/handler"
	"github.com/vatsal278/msgbroker/internal/model"
	parser "github.com/vatsal278/msgbroker/internal/pkg/parser"
	RSA "github.com/vatsal278/msgbroker/pkg/Rsa"
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
		//PublicKey:="24049099257750543577110691145392378535433506714382016684543650515257577861982800257875667193131828538020295509676212580314884562301451353127727792254976742096165542847693603938683809531049339684700695684621360161384777600216173105205424925494517556647992162280299092209488690220952214623117409133793117229918719355861564257788861574556262434611828388560502280837107250951428122252993587372819460884050317262549673940088624322772988739161204172820022948311146833791624240432914774891660328641595103012681073823066179394147279398819054276953900486280012050787242661250538787611882539749025308645087482171961775641573583 65537"

		for _, v := range m.messageBroker.SubM[updates.Publisher.Channel] {
			go func(v model.Subscriber) {
				PublicKey := v.CallBack.PublicKey
				PubKey := RSA.PEMStrAsKey(PublicKey)
				updates.Update = RSA.RSA_OAEP_Encrypt(updates.Update, *PubKey)
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
