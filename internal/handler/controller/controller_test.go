package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vatsal278/msgbroker/internal/constants"
	controllerInterface "github.com/vatsal278/msgbroker/internal/handler"
	"github.com/vatsal278/msgbroker/internal/model"
)

type tempStruct struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func TestRegisterPublisher(t *testing.T) {
	var publisher = model.Publisher{
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
				var y model.Publisher = reqbody.(model.Publisher)
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
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
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
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
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

	var callback = model.CallBack{
		HttpMethod:  "GET",
		CallbackUrl: "http://localhost:8083/pong",
	}

	var subscriber = model.Subscriber{
		CallBack: callback,
		Channel:  "c4",
	}

	type TempSubscriber struct {
		CallBack model.CallBack
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
				var y model.Subscriber = reqbody.(model.Subscriber)
				t.Log(y)
				//m := x.messageBroker.PubM[publisher.Channel]

				for {
					m := x.messageBroker.SubM[subscriber.Channel]
					if len(m) == 1 {
						break
					}
				}
				m := x.messageBroker.SubM[subscriber.Channel]
				if len(m) != 1 {
					t.Errorf("Want: %v, Got: %v", "1", len(m))
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
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
				if len(m) != 0 {
					t.Errorf("Want: %v, Got: %v", "not ok", len(m))

				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
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
func DummyRegister(url string, method string, t *testing.T, i controllerInterface.IController) {
	var callback = model.CallBack{
		HttpMethod:  method,
		CallbackUrl: url,
	}
	var subscriber = model.Subscriber{
		CallBack: callback,
		Channel:  "c4",
	}
	var m *models = i.(*models)
	//var y TempSubscriber = reqbody.(TempSubscriber)
	subs := m.messageBroker.SubM[subscriber.Channel]

	subs = append(subs, subscriber)
	t.Log(subs)
	t.Logf("subscriber added %+v", subscriber)
	m.messageBroker.SubM[subscriber.Channel] = subs
	t.Log(m.messageBroker.SubM[subscriber.Channel])
}

func Testutility(c *T) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		defer c.wg.Done()
		defer c.t.Log("HIT")
		x, err := ioutil.ReadAll(r.Body)
		if err != nil {
			c.t.Log(err.Error())
		}
		if !reflect.DeepEqual(string(x), "Hello World") {
			c.t.Errorf("Want: %v, Got: %v", "Hello World", string(x))
		}
	}).Methods(http.MethodPost)
	http.Handle("/", router)
	fmt.Println("Connected to Test Server")

	return router
}
func testClient(c *T) {
	//expected := "dummy data"

	x := Testutility(c)
	svr := httptest.NewServer(x)
	url := svr.URL + "/ping"
	c.t.Log(url)
	c.srv = svr
	//defer svr.Close()
	DummyRegister(url, "POST", c.t, c.i)

}

type T struct {
	srv *httptest.Server
	t   *testing.T
	i   controllerInterface.IController
	wg  *sync.WaitGroup
}

func TestPublishMessage(t *testing.T) {

	var publisher = model.Publisher{
		Name:    "publisher1",
		Channel: "c4",
	}
	type TempPublisher struct {
		Name    interface{} `validate:"required"`
		Channel string      `form:"channel" json:"channel" validate:"required"`
	}

	var updates = model.Updates{
		Publisher: publisher,
		Update:    "Hello World",
	}

	type TempUpdates struct {
		Publisher TempPublisher `form:"publisher" json:"publisher" validate:"required"`
		Update    int           `form:"update" json:"update" validate:"required"`
	}
	var tempPublisher = TempPublisher{
		Name:    1,
		Channel: "c4",
	}
	var dummy = TempUpdates{
		Publisher: tempPublisher,
		Update:    1,
	}

	tStruct := &T{
		t:  t,
		wg: &sync.WaitGroup{},
	}

	tests := []struct {
		name             string
		requestBody      interface{}
		expectedResponse tempStruct
		setupFunc        func(controllerInterface.IController)
		validateFunc     func(*httptest.ResponseRecorder)
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: updates,
			setupFunc: func(i controllerInterface.IController) {
				tStruct.i = i
				tStruct.wg.Add(1)
				testClient(tStruct)

				var x *models = i.(*models)
				m, ok := x.messageBroker.PubM[publisher.Channel]
				if !ok {
					m = make(map[string]struct{})
					m[publisher.Name] = struct{}{}
				}
				x.messageBroker.PubM[publisher.Channel] = m
				t.Log(x.messageBroker.SubM)
			},
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = tempStruct{
					Status:  http.StatusOK,
					Message: "notified all subscriber",
					Data:    nil,
				}
				tStruct.wg.Wait()
				tStruct.srv.Close()
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				if w.Code != tempstruct.Status {
					t.Errorf("Want: %v, Got: %v", tempstruct.Status, w.Code)
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
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
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
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = tempStruct{
					Status:  http.StatusNotFound,
					Message: "No publisher found with the specified name for specified channel",
					Data:    nil,
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				if w.Code != tempstruct.Status {
					t.Errorf("Want: %v, Got: %v", tempstruct.Status, w.Code)
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
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
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
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = tempStruct{
					Status:  http.StatusBadRequest,
					Message: constants.IncompleteData,
					Data:    nil,
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				if w.Code != tempstruct.Status {
					t.Errorf("Want: %v, Got: %v", tempstruct.Status, w.Code)
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
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
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
			Publish := i.PublishMessage()
			tt.setupFunc(i)
			w := httptest.NewRecorder()
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/publish", bytes.NewBuffer(jsonValue))

			Publish(w, r)

			if tt.validateFunc != nil {
				tt.validateFunc(w)
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

			validateFunc: func(w *httptest.ResponseRecorder, r *http.Request) {
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				expectedResponse := Response{
					Status:  http.StatusNotFound,
					Message: "no route found",
					Data:    nil,
				}
				if w.Code != expectedResponse.Status {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, w.Code)
				}
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
