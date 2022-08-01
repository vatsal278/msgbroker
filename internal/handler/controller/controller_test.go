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
	/*var callback = model.CallBack{
		HttpMethod:  "GET",
		CallbackUrl: "http://localhost:8083/pong",
	}
	var subscriber = model.Subscriber{
		CallBack: callback,
		Channel:  "c4",
	}*/
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
		expectedResponse temp_struct
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
				/*w := httptest.NewRecorder()
				jsonValue, _ := json.Marshal(subscriber)
				r := httptest.NewRequest("POST", "/register/subscriber", bytes.NewBuffer(jsonValue))
				RegisterSub := i.RegisterSubscriber()
				RegisterSub(w, r)*/

			},
			expectedResponse: temp_struct{
				Status:  http.StatusNotFound,
				Message: "No publisher found with the specified name for specified channel",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE::Parse error",
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
			expectedResponse: temp_struct{
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

func TestNoArticleFound(t *testing.T) {

	w := httptest.NewRecorder()
	//errorCase := int(tt.ErrorCase)
	i := NewController()
	NorouteController := i.NoRouteFound()
	r := httptest.NewRequest("POST", "/a", nil)
	NorouteController(w, r)
	response_body, error := ioutil.ReadAll(w.Body)
	if error != nil {
		t.Error(error.Error())
	}
	var response model.Response
	err := json.Unmarshal(response_body, &response)
	if err != nil {
		t.Error(error.Error())
	}
}
