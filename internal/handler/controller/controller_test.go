package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vatsal278/msgbroker/internal/constants"
	controllerInterface "github.com/vatsal278/msgbroker/internal/handler"
)

type tempStruct struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func TestRegisterPublisher(t *testing.T) {
	type Publisher struct {
		Name    string `form:"name" json:"name" validate:"required"`
		Channel string `form:"channel" json:"channel" validate:"required"`
	}
	var publisher = Publisher{
		Name:    "publisher1",
		Channel: "c4",
	}
	type tempPublisher struct {
		Name    interface{} `validate:"required"`
		Channel string      `form:"channel" json:"channel" validate:"required"`
	}
	var dummy = tempPublisher{
		Name:    1,
		Channel: "c4",
	}
	tests := []struct {
		name              string
		requestBody       interface{}
		ValidateFunc      func(*httptest.ResponseRecorder, controllerInterface.IController, interface{})
		expected_response tempStruct
	}{
		{
			name:        "Success:: Register Publisher",
			requestBody: publisher,
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y Publisher = reqbody.(Publisher)
				t.Log(y)

				m, ok := x.messageBroker.PubM[publisher.Channel]
				if !ok {
					t.Errorf("Want: %v, Got: %v", "publisher map", ok)
				}
				_, ok = m[publisher.Name]
				if !ok {
					t.Errorf("Want: %v, Got: %v", "publisher map", ok)
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v", nil)
				}
				if w.Code != http.StatusCreated {
					t.Errorf("Want: %v, Got: %v", http.StatusCreated, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response tempStruct
				err := json.Unmarshal(responseBody, &response)
				expectedResponse := tempStruct{
					Status:  http.StatusCreated,
					Message: "Successfully Registered as publisher to the channel",
					Data:    nil,
				}
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, response)
				}
			},
		},
		{
			name:        "FAILURE:: Register Publisher:Incorrect Input Details",
			requestBody: dummy,
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y tempPublisher = reqbody.(tempPublisher)
				t.Log(y)

				m, ok := x.messageBroker.PubM[dummy.Channel]
				if ok {
					t.Errorf("Want: %v, Got: %v", "not ok", ok)
				}
				_, ok = m[publisher.Name]
				if ok {
					t.Errorf("Want: %v, Got: %v", "not ok", ok)
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v", nil)
				}
				if w.Code != http.StatusBadRequest {
					t.Errorf("Want: %v, Got: %v", http.StatusBadRequest, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response tempStruct
				err := json.Unmarshal(responseBody, &response)
				expectedResponse := tempStruct{
					Status:  http.StatusBadRequest,
					Message: constants.IncompleteData,
					Data:    nil,
				}
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			t.Log(w.Code)
			reqbody := tt.requestBody
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			i := NewController()
			RegisterPub := i.RegisterPublisher()
			RegisterPub(w, r)
			tt.ValidateFunc(w, i, reqbody)
		})
	}
}

func TestRegisterSubscriber(t *testing.T) {
	type CallBack struct {
		HttpMethod  string `validate:"required"`
		CallbackUrl string `validate:"required"`
	}
	var callback = CallBack{
		HttpMethod:  "GET",
		CallbackUrl: "http://localhost:8083/pong",
	}
	type Subscriber struct {
		CallBack CallBack
		Channel  string `form:"channel" json:"channel" validate:"required"`
	}
	var subscriber = Subscriber{
		CallBack: callback,
		Channel:  "c4",
	}

	type TempSubscriber struct {
		CallBack CallBack
		Channel  int `form:"channel" json:"channel" validate:"required"`
	}
	var dummy = TempSubscriber{
		CallBack: callback,
		Channel:  1,
	}
	tests := []struct {
		name              string
		expected_response tempStruct
		requestBody       interface{}
		ValidateFunc      func(*httptest.ResponseRecorder, controllerInterface.IController, interface{})
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: subscriber,
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y Subscriber = reqbody.(Subscriber)
				t.Log(y)
				//m := x.messageBroker.PubM[publisher.Channel]

				for {
					m := x.messageBroker.SubM[subscriber.Channel]
					if len(m) == 1 {
						break
					}
				}
				m := x.messageBroker.SubM[subscriber.Channel]
				if len(m) == 0 {
					t.Errorf("Want: %v, Got: %v", "1", len(m))
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v", nil)
				}
				expectedResponse := tempStruct{
					Status:  http.StatusCreated,
					Message: "Successfully Registered as Subscriber to the channel",
					Data:    nil,
				}

				if w.Code != expectedResponse.Status {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response tempStruct
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, response)
				}
			},
		},
		{
			name:        "FAILURE:: Register subscriber::Incorrect Input Details",
			requestBody: dummy,
			expected_response: tempStruct{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				//var y TempSubscriber = reqbody.(TempSubscriber)
				m := x.messageBroker.SubM[subscriber.Channel]
				for {
					m := x.messageBroker.SubM[subscriber.Channel]
					if len(m) == 0 {
						break
					}
				}
				if len(m) != 0 {
					t.Errorf("Want: %v, Got: %v", "not ok", len(m))

				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v", nil)
				}
				expectedResponse := tempStruct{
					Status:  http.StatusBadRequest,
					Message: constants.IncompleteData,
					Data:    nil,
				}

				if w.Code != expectedResponse.Status {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response tempStruct
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, response)
				}
			},
		},
	}
	//creating separate validate func
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			t.Log(w.Code)
			reqBody := tt.requestBody
			jsonValue, _ := json.Marshal(reqBody)
			r := httptest.NewRequest("POST", "/register/subscriber", bytes.NewBuffer(jsonValue))
			i := NewController()
			RegisterSub := i.RegisterSubscriber()
			RegisterSub(w, r)
			tt.ValidateFunc(w, i, reqBody)
		})
	}
}

func TestPublishMessage(t *testing.T) {

	type Publisher struct {
		Name    string `form:"name" json:"name" validate:"required"`
		Channel string `form:"channel" json:"channel" validate:"required"`
	}
	var publisher = Publisher{
		Name:    "publisher1",
		Channel: "c4",
	}
	type TempPublisher struct {
		Name    interface{} `validate:"required"`
		Channel string      `form:"channel" json:"channel" validate:"required"`
	}
	var tempPublisher = TempPublisher{
		Name:    1,
		Channel: "c4",
	}
	type Updates struct {
		Publisher Publisher `form:"publisher" json:"publisher" validate:"required"`
		Update    string    `form:"update" json:"update" validate:"required"`
	}
	var updates = Updates{
		Publisher: publisher,
		Update:    "Hello World",
	}

	type TempUpdates struct {
		Publisher TempPublisher `form:"publisher" json:"publisher" validate:"required"`
		Update    int           `form:"update" json:"update" validate:"required"`
	}
	var dummy = TempUpdates{
		Publisher: tempPublisher,
		Update:    1,
	}

	tests := []struct {
		name             string
		requestBody      interface{}
		expectedResponse tempStruct
		setupFunc        func(controllerInterface.IController)
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: updates,
			setupFunc: func(i controllerInterface.IController) {
				var x *models = i.(*models)
				m, ok := x.messageBroker.PubM[publisher.Channel]
				if !ok {
					m = make(map[string]struct{})
					m[publisher.Name] = struct{}{}
				}
				x.messageBroker.PubM[publisher.Channel] = m
			},
			expectedResponse: tempStruct{
				Status:  http.StatusOK,
				Message: "notified all subscriber",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE::Publish Message::No Publisher Found",
			requestBody: updates,
			setupFunc: func(i controllerInterface.IController) {
			},
			expectedResponse: tempStruct{
				Status:  http.StatusNotFound,
				Message: "No publisher found with the specified name for specified channel",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE::Publish Message::Incorrect input details",
			requestBody: dummy,
			setupFunc: func(i controllerInterface.IController) {
				var x *models = i.(*models)
				m, ok := x.messageBroker.PubM[publisher.Channel]
				if !ok {
					m = make(map[string]struct{})
					m[publisher.Name] = struct{}{}
				}
				x.messageBroker.PubM[publisher.Channel] = m
			},
			expectedResponse: tempStruct{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//golang http test server
			i := NewController()

			tt.setupFunc(i)
			w := httptest.NewRecorder()
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/publish", bytes.NewBuffer(jsonValue))
			Publish := i.PublishMessage()
			Publish(w, r)
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Want: Content Type as %v", nil)
			}
			if w.Code != tt.expectedResponse.Status {
				t.Errorf("Want: %v, Got: %v", tt.expectedResponse.Status, w.Code)
			}
			responseBody, error := ioutil.ReadAll(w.Body)
			if error != nil {
				t.Error(error.Error())
			}
			var response tempStruct
			err := json.Unmarshal(responseBody, &response)
			if err != nil {
				t.Error(error.Error())
			}
			t.Log(response)
			if !reflect.DeepEqual(response, tt.expectedResponse) {
				t.Errorf("Want: %v, Got: %v", tt.expectedResponse, response)
			}
		})
	}
}

func TestNoRouteFound(t *testing.T) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    interface{}
	}
	tests := []struct {
		name             string
		requestBody      interface{}
		expectedResponse Response
		validateFunc     func(*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "Success:: NoRouteFound",

			expectedResponse: Response{
				Status:  http.StatusNotFound,
				Message: "no route found",
				Data:    nil,
			},
			validateFunc: func(w *httptest.ResponseRecorder, r *http.Request) {
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response Response
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				t.Log(response)
				expectedResponse := Response{
					Status:  http.StatusNotFound,
					Message: "No Route Found",
					Data:    nil,
				}
				if !reflect.DeepEqual(response, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			i := NewController()
			NorouteController := i.NoRouteFound()
			r := httptest.NewRequest("POST", "/a", nil)
			NorouteController(w, r)
			tt.validateFunc(w, r)
		})
	}
}
