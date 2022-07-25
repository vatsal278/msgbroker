package article_controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/model"

	"github.com/go-playground/validator"
)

var SubscriberMap = map[string][]model.Subscriber{}
var PublisherMap = map[string]map[string]struct{}{}

var MessageBroker = model.MessageBroker{
	SubM: SubscriberMap,
	PubM: PublisherMap,
}

type IController interface {
	RegisterPublisher() func(w http.ResponseWriter, r *http.Request)
	RegisterSubscriber() func(w http.ResponseWriter, r *http.Request)
	PublishMessage() func(w http.ResponseWriter, r *http.Request)
}

func RegisterPublisher() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var publisher model.Publisher
		w.Header().Set("Content-Type", "application/json")
		//Read body of the request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		// parse json encoded data into structure
		err = json.Unmarshal(body, &publisher)

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		//Create separete pkg directory
		validate := validator.New()
		errs := validate.Struct(publisher)
		if errs != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusBadRequest,
				constants.Incomplete_Data, nil))
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		x, ok := MessageBroker.PubM[publisher.Channel]

		//

		//

		//
		if !ok {
			x = make(map[string]struct{})
			x[publisher.Channel] = publisher
			MessageBroker.PubM[publisher.Channel] = x
		}
		MessageBroker.PubM[publisher.Channel] = x

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Successfully Registered as publisher to the channel", nil))
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
		//Read body of the request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		// parse json encoded data into structure
		err = json.Unmarshal(body, &subscriber)

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		validate := validator.New()
		errs := validate.Struct(subscriber)
		if errs != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusBadRequest,
				constants.Incomplete_Data, nil))
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func(s model.Subscriber) {
			defer wg.Done()
			MessageBroker.Lock()
			defer MessageBroker.Unlock()
			subs := MessageBroker.SubM[s.Channel]

			for _, v := range subs {
				if reflect.DeepEqual(v, s) {
					w.WriteHeader(http.StatusCreated)
					err = json.NewEncoder(w).Encode(Response_Writer(http.StatusOK, "Subscriber already exists", nil))
					if err != nil {
						log.Println(err.Error())
					}
					log.Print("Subscriber already exists")
					return
				}
			}
			subs = append(subs, s)
			MessageBroker.SubM[s.Channel] = subs

			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Successfully Subscribed to the channel", nil))
			if err != nil {
				log.Println(err.Error())
			}
			log.Print("Successfully Subscribed to the channel")
		}(subscriber)
		wg.Wait()
	}
}

func PublishMessage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		w.Header().Set("Content-Type", "application/json")
		//Read body of the request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		// parse json encoded data into structure
		err = json.Unmarshal(body, &updates)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		validate := validator.New()
		errs := validate.Struct(updates)
		if errs != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusBadRequest,
				constants.Incomplete_Data, nil))
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)

		pubm := MessageBroker.PubM[updates.Publisher.Channel]
		_, ok := pubm[updates.Publisher.Name]
		if !ok {
			err = json.NewEncoder(w).Encode(Response_Writer(http.StatusNotFound, "No publisher found with the specified name for specified channel", nil))
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

		//ChannelUpdates = append(ChannelUpdates, updates

		for _, v := range MessageBroker.SubM[updates.Publisher.Channel] {
			go func(v model.Subscriber) {
				MessageBroker.Lock()
				defer MessageBroker.Unlock()
				log.Print("sending notification")
				//Call another route to notify publisher
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
				//wg.Wait()
				//log.Print("sent updates to the channel")
			}(v)

		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Sending notification", nil))
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func NotifySubscriber() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates model.Updates
		w.Header().Set("Content-Type", "application/json")
		//Read body of the request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		// parse json encoded data into structure
		err = json.Unmarshal(body, &updates)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(Response_Writer(http.StatusInternalServerError, constants.Parse_Err, nil))
			if err != nil {
				log.Println(err.Error())
			}

			return
		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "received updates on the channel", nil))
		if err != nil {
			log.Println(err.Error())
		}
		log.Print("Received updates on the channel")
	}
}

func Response_Writer(status int, msg string, data interface{}) model.Response {
	var response model.Response
	response.Status = status
	response.Message = msg
	response.Data = data
	return response
}
