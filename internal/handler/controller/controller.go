package article_controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"

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
		x := MessageBroker.PubM[publisher.Channel]
		MessageBroker.PubM[publisher.Channel] = x

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Successfully Subscribed to the channel", nil))
		if err != nil {
			log.Println(err.Error())
		}
		log.Print("Successfully Subscribed to the channel")
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
		go func(p model.Publisher) {
			defer wg.Done()
			MessageBroker.Lock()
			defer MessageBroker.Unlock()
			pubm := MessageBroker.PubM[p.Channel]
			pubs := pubm[updates.Publisher.Name]

			for _, v := range pubm {
				if reflect.DeepEqual(v, pubs) {
					err = json.NewEncoder(w).Encode(Response_Writer(http.StatusNotFound, "No publisher found with the specified name for specified channel", nil))
					if err != nil {
						log.Println(err.Error())
					}
					return
				}
			}
		}(updates.Publisher)
		wg.Wait()
		//ChannelUpdates = append(ChannelUpdates, updates)
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Successfully sent updates to the channel", nil))
		if err != nil {
			log.Println(err.Error())
		}
		log.Print("Successfully sent updates to the channel")
		for _, v := range MessageBroker.SubM[updates.Publisher.Channel] {
			go func() {
				defer wg.Done()
				MessageBroker.Lock()
				defer MessageBroker.Unlock()
				log.Print("sending notification")
				//Call another route to notify publisher
				reqBody, err := json.Marshal(updates.Update)
				if err != nil {
					log.Println(err.Error())
				}
				//timeout := time.Duration(2 * time.Second)
				client := http.DefaultClient
				method := v.CallBack.HttpMethod
				url := v.CallBack.CallbackUrl
				request, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Println(err.Error())
				}
				log.Printf("%+v \n", *request)
				client.Do(request)
				wg.Wait()
			}()
			wg.Wait()

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
