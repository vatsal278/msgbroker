package article_controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/model"

	"github.com/go-playground/validator"
)

var SubscriberList []model.Subscriber
var PublisherList []model.Publisher
var ChannelUpdates []model.Updates

type IController interface {
	RegisterPublisher() func(w http.ResponseWriter, r *http.Request)
	RegisterSubscriber() func(w http.ResponseWriter, r *http.Request)
	PublishMessage() func(w http.ResponseWriter, r *http.Request)
}
type TController interface {
	NotifySubscriber() func(w http.ResponseWriter, r *http.Request)
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

		go func(publisher model.Publisher) {
			for _, pub := range PublisherList {

				if reflect.DeepEqual(model.Publisher(publisher), pub) {
					return
				}
			}
			PublisherList = append(PublisherList, publisher)
		}(publisher)

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
		count := 0
		for _, sub := range SubscriberList {

			if sub == subscriber {
				count += 1
			}
		}
		if count == 0 {
			SubscriberList = append(SubscriberList, subscriber)
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Successfully Subscribed to the channel", nil))
			if err != nil {
				log.Println(err.Error())
			}
			log.Print("Successfully Subscribed to the channel")
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(Response_Writer(http.StatusOK, "Subscriber already exists", nil))
			if err != nil {
				log.Println(err.Error())
			}
			log.Print("Successfully Subscribed to the channel")
			return
		}
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
		count := 0
		for _, pub := range PublisherList {

			if pub.Publisher == updates.Publisher.Publisher {
				count += 1
			}
		}
		if count == 0 {
			err = json.NewEncoder(w).Encode(Response_Writer(http.StatusNotFound, "No publisher found with the specified name for specified channel", nil))
			if err != nil {
				log.Println(err.Error())
			}
			return
		} else {
			ChannelUpdates = append(ChannelUpdates, updates)
			w.WriteHeader(http.StatusCreated)

			err = json.NewEncoder(w).Encode(Response_Writer(http.StatusCreated, "Successfully sent updates to the channel", nil))
			if err != nil {
				log.Println(err.Error())
			}
			log.Print("Successfully sent updates to the channel")
		}
		for _, v := range SubscriberList {
			if v.Channel == "c1" {
				log.Print("sending notification")
				//do this
			}
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
