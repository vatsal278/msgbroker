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
	"github.com/vatsal278/msgbroker/internal/model"
)

type temp_struct struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func TestRegisterPublisher(t *testing.T) {
	var publisher = model.Publisher{
		Name:    "publisher1",
		Channel: "c4",
	}
	var dummy = model.TempPublisher{
		Name:    1,
		Channel: "c4",
	}
	tests := []struct {
		name              string
		requestBody       interface{}
		ErrorCase         bool
		expected_response temp_struct
	}{
		{
			name:        "Success:: Register Publisher",
			requestBody: publisher,
			expected_response: temp_struct{
				Status:  http.StatusCreated,
				Message: "Successfully Registered as publisher to the channel",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE:: Register Publisher",
			requestBody: dummy,
			expected_response: temp_struct{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			t.Log(w.Code)
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			i := NewController()
			RegisterPub := i.RegisterPublisher()
			RegisterPub(w, r)
			response_body, error := ioutil.ReadAll(w.Body)
			if error != nil {
				t.Error(error.Error())
			}
			var response temp_struct
			err := json.Unmarshal(response_body, &response)
			if err != nil {
				t.Error(error.Error())
			}
			if !reflect.DeepEqual(response, tt.expected_response) {
				t.Errorf("Want: %v, Got: %v", tt.expected_response, response)
			}
		})
	}
}

func TestRegisterSubscriber(t *testing.T) {
	var callback = model.CallBack{
		HttpMethod:  "GET",
		CallbackUrl: "http://localhost:8083/pong",
	}
	var subscriber = model.Subscriber{
		CallBack: callback,
		Channel:  "c4",
	}
	var dummy = model.TempSubscriber{
		CallBack: callback,
		Channel:  1,
	}
	tests := []struct {
		name              string
		expected_response temp_struct
		requestBody       interface{}
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: subscriber,
			expected_response: temp_struct{
				Status:  http.StatusCreated,
				Message: "Successfully Registered as Subscriber to the channel",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE:: Register Publisher",
			requestBody: dummy,
			expected_response: temp_struct{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			t.Log(w.Code)
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/subscriber", bytes.NewBuffer(jsonValue))
			i := NewController()
			RegisterSub := i.RegisterSubscriber()
			RegisterSub(w, r)
			response_body, error := ioutil.ReadAll(w.Body)
			if error != nil {
				t.Error(error.Error())
			}
			var response temp_struct
			err := json.Unmarshal(response_body, &response)
			if err != nil {
				t.Error(error.Error())
			}
			if !reflect.DeepEqual(response, tt.expected_response) {
				t.Errorf("Want: %v, Got: %v", tt.expected_response, response)
			}
		})
	}
}

func TestPublishMessage(t *testing.T) {
	var callback = model.CallBack{
		HttpMethod:  "GET",
		CallbackUrl: "http://localhost:8083/pong",
	}
	var subscriber = model.Subscriber{
		CallBack: callback,
		Channel:  "c4",
	}
	var publisher = model.Publisher{
		Name:    "publisher1",
		Channel: "c4",
	}
	var tempPublisher = model.TempPublisher{
		Name:    1,
		Channel: "c4",
	}
	var updates = model.Updates{
		Publisher: publisher,
		Update:    "Hello World",
	}
	var dummy = model.TempUpdates{
		Publisher: tempPublisher,
		Update:    1,
	}

	tests := []struct {
		name             string
		requestBody      interface{}
		expectedResponse temp_struct
		setupFunc        func(controllerInterface.IController)
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: updates,
			setupFunc: func(i controllerInterface.IController) {
				w := httptest.NewRecorder()
				jsonValue, _ := json.Marshal(publisher)
				r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
				RegisterPub := i.RegisterPublisher()
				RegisterPub(w, r)
				w = httptest.NewRecorder()
				jsonValue, _ = json.Marshal(subscriber)
				r = httptest.NewRequest("POST", "/register/subscriber", bytes.NewBuffer(jsonValue))
				RegisterSub := i.RegisterSubscriber()
				RegisterSub(w, r)
			},
			expectedResponse: temp_struct{
				Status:  http.StatusOK,
				Message: "notified all subscriber",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE:: no publisher found",
			requestBody: updates,
			setupFunc: func(i controllerInterface.IController) {
				w := httptest.NewRecorder()
				jsonValue, _ := json.Marshal(subscriber)
				r := httptest.NewRequest("POST", "/register/subscriber", bytes.NewBuffer(jsonValue))
				RegisterSub := i.RegisterSubscriber()
				RegisterSub(w, r)
			},
			expectedResponse: temp_struct{
				Status:  http.StatusNotFound,
				Message: "No publisher found with the specified name for specified channel",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE::P",
			requestBody: dummy,
			setupFunc: func(i controllerInterface.IController) {
				w := httptest.NewRecorder()
				jsonValue, _ := json.Marshal(publisher)
				r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
				RegisterPub := i.RegisterPublisher()
				RegisterPub(w, r)
				w = httptest.NewRecorder()
				jsonValue, _ = json.Marshal(subscriber)
				r = httptest.NewRequest("POST", "/register/subscriber", bytes.NewBuffer(jsonValue))
				RegisterSub := i.RegisterSubscriber()
				RegisterSub(w, r)
			},
			expectedResponse: temp_struct{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewController()
			tt.setupFunc(i)
			w := httptest.NewRecorder()
			jsonValue, _ := json.Marshal(updates)
			r := httptest.NewRequest("POST", "/publish", bytes.NewBuffer(jsonValue))
			Publish := i.PublishMessage()
			Publish(w, r)
			response_body, error := ioutil.ReadAll(w.Body)
			if error != nil {
				t.Error(error.Error())
			}
			var response temp_struct
			err := json.Unmarshal(response_body, &response)
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
