package controller

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"github.com/vatsal278/msgbroker/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/handler"
)

func TestRegisterPublisher(t *testing.T) {
	type tempPublisher struct {
		Channel interface{} `form:"channel" json:"channel" validate:"required"`
	}

	tests := []struct {
		name             string
		requestBody      interface{}
		ValidateFunc     func(*httptest.ResponseRecorder, controllerInterface.IController, interface{})
		expectedResponse model.Response
	}{
		{
			name: "Success:: Register Publisher",
			requestBody: model.Publisher{
				Channel: "c4",
			},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y model.Publisher = reqbody.(model.Publisher)
				t.Log(y)

				m, ok := x.messageBroker.PubM["c4"]
				if !ok {
					t.Errorf("Want: %v, Got: %v", "publisher map", ok)
				}
				if len(m) == 0 {
					t.Errorf("Want: %v, Got: %v", "publisher map", m)
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				expectedResponse := model.Response{
					Status:  http.StatusCreated,
					Message: constants.PublisherRegistration,
					Data: map[string]interface{}{
						"id": "57409864-9a6e-4595-a8fa-ac7e3a61da74",
					},
				}
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response.Status, expectedResponse.Status) {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, response.Status)
				}
				if !reflect.DeepEqual(response.Message, expectedResponse.Message) {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Message, response.Message)
				}
				b := response.Data
				var a map[string]interface{} = b.(map[string]interface{})
				var z string = a["id"].(string)
				_, err = uuid.Parse(z)

				if err != nil {
					t.Error(err.Error())
				}
				_, ok = m[z]
				if !ok {
					t.Errorf("Want: %v, Got: %v", "publisher map", ok)
				}
			},
		},
		{
			name: "FAILURE:: Register Publisher:Incorrect Input Details",
			requestBody: tempPublisher{
				Channel: 1,
			},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y tempPublisher = reqbody.(tempPublisher)
				t.Log(y)
				m := x.messageBroker.PubM["c4"]
				if len(m) != 0 {
					t.Errorf("Want: %v, Got: %v", nil, m)
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				expectedResponse := model.Response{
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

	type TempSubscriber struct {
		CallBack model.CallBack
		Channel  int `form:"channel" json:"channel" validate:"required"`
	}

	tests := []struct {
		name             string
		expectedResponse model.Response
		requestBody      interface{}
		ValidateFunc     func(*httptest.ResponseRecorder, controllerInterface.IController, interface{})
	}{
		{
			name: "Success:: Register Subscriber::With Encryption",
			requestBody: model.Subscriber{
				CallBack: model.CallBack{
					HttpMethod:  "GET",
					CallbackUrl: "http://localhost:8083/pong",
					PublicKey:   "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJDZ0tDQVFFQXZBWmZxM1lvVzdUTzBGYmJHMWxxRVBxNHQ4bGc5cTdla0NYMXJIVjVNNTdobmdyNlF1L3MKTnp0QXkzTmh1TG4xSm5PSVN5bzRXc29MMDRKWFI5WXI5UXVtZW1EdGVreWpOd2toQkFWM0xBN3BORjV3c2ZaSwpFbC9jY2U5aGZxRWtOcERtNUFFZklnRW5UZXdTMml5cGRCQm1pVmI5VzNzZFdUWHEwenNKY1pqb29obXZPNkN1CngyY01NOW1EeFQ4VXBYM2gweE1WNTBVd050TzRVbS9aWnFPeENqdFdhNE1STE16NTNMTG9lUm9UOE1tZEdlV1UKYTdHMitKU0c5K3V1MVJIVkYrelZGaEx2emtoM3dLTGdVdU1DcW0rL1U0Y3B3TDUxZU9TYVZNYUhjU1NiRXZCUgp0d0lZdHRHR3NDVC9mTEdyVXdjZm8xZ0xKaVNjU2taN1B3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
				},
				Channel: "c4",
			},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y model.Subscriber = reqbody.(model.Subscriber)
				t.Log(y)
				//m := x.messageBroker.PubM[publisher.Channel]

				for {
					m := x.messageBroker.SubM["c4"]
					if len(m) == 1 {
						break
					}
				}
				m := x.messageBroker.SubM["c4"]
				if len(m) != 1 {
					t.Errorf("Want: %v, Got: %v", "1", len(m))
				}

				if !reflect.DeepEqual(m[0], y) {
					t.Errorf("Want: %v, Got: %v", y, m[0])
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				expectedResponse := model.Response{
					Status:  http.StatusCreated,
					Message: constants.SubscriberRegistration,
					Data:    nil,
				}

				if w.Code != expectedResponse.Status {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response model.Response
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
			name: "Success:: Register Subscriber :: Without Encryption",
			requestBody: model.Subscriber{
				CallBack: model.CallBack{
					HttpMethod:  "GET",
					CallbackUrl: "http://localhost:8083/pong",
				},
				Channel: "c4",
			},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				var y model.Subscriber = reqbody.(model.Subscriber)
				t.Log(y)
				//m := x.messageBroker.PubM[publisher.Channel]

				for {
					m := x.messageBroker.SubM[y.Channel]
					if len(m) == 1 {
						break
					}
				}
				m := x.messageBroker.SubM[y.Channel]
				if len(m) != 1 {
					t.Errorf("Want: %v, Got: %v", "1", len(m))
				}
				if !reflect.DeepEqual(m[0], y) {
					t.Errorf("Want: %v, Got: %v", y, m[0])
				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				expectedResponse := model.Response{
					Status:  http.StatusCreated,
					Message: constants.SubscriberRegistration,
					Data:    nil,
				}

				if w.Code != expectedResponse.Status {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response model.Response
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
			name: "FAILURE:: Register subscriber::Incorrect Input Details",
			requestBody: TempSubscriber{
				CallBack: model.CallBack{
					HttpMethod:  "GET",
					CallbackUrl: "http://localhost:8083/pong",
					PublicKey:   "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJDZ0tDQVFFQXZBWmZxM1lvVzdUTzBGYmJHMWxxRVBxNHQ4bGc5cTdla0NYMXJIVjVNNTdobmdyNlF1L3MKTnp0QXkzTmh1TG4xSm5PSVN5bzRXc29MMDRKWFI5WXI5UXVtZW1EdGVreWpOd2toQkFWM0xBN3BORjV3c2ZaSwpFbC9jY2U5aGZxRWtOcERtNUFFZklnRW5UZXdTMml5cGRCQm1pVmI5VzNzZFdUWHEwenNKY1pqb29obXZPNkN1CngyY01NOW1EeFQ4VXBYM2gweE1WNTBVd050TzRVbS9aWnFPeENqdFdhNE1STE16NTNMTG9lUm9UOE1tZEdlV1UKYTdHMitKU0c5K3V1MVJIVkYrelZGaEx2emtoM3dLTGdVdU1DcW0rL1U0Y3B3TDUxZU9TYVZNYUhjU1NiRXZCUgp0d0lZdHRHR3NDVC9mTEdyVXdjZm8xZ0xKaVNjU2taN1B3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
				},
				Channel: 1,
			},
			expectedResponse: model.Response{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {
				var x *models = i.(*models)
				//var y TempSubscriber = reqbody.(TempSubscriber)
				m := x.messageBroker.SubM["c4"]
				if len(m) != 0 {
					t.Errorf("Want: %v, Got: %v", "not ok", len(m))

				}
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				expectedResponse := model.Response{
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
				var response model.Response
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

func testClient(c *TestServer, encrypted bool) {
	//expected := "dummy data"
	var privateKey *rsa.PrivateKey
	var err error

	if encrypted {
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)

		if err != nil {
			c.t.Log(err.Error())
		}

	}
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		defer c.wg.Done()
		defer c.t.Log("HIT")
		x, err := ioutil.ReadAll(r.Body)
		if err != nil {
			c.t.Log(err.Error())
			return
		}
		res := string(x)

		if privateKey != nil {
			res, err = crypt.RsaOaepDecrypt(string(x), *privateKey)

		}

		if !reflect.DeepEqual(res, "Hello World") {
			c.t.Errorf("Want: %v, Got: %v", "Hello World", res)
			return
		}
		c.t.Log(res)
	}).Methods(http.MethodPost)
	svr := httptest.NewServer(router)
	var url = svr.URL + "/ping"
	c.srv = svr
	c.t.Log(privateKey)
	var pubKey = ""
	if privateKey != nil {
		publicKey := privateKey.PublicKey
		pubKey = crypt.KeyAsPEMStr(&publicKey)
		c.t.Log(pubKey)
	}

	var subscriber = model.Subscriber{
		CallBack: model.CallBack{
			HttpMethod:  "POST",
			CallbackUrl: url,
			PublicKey:   pubKey,
		},
		Channel: "c4",
	}
	var m *models = c.i.(*models)
	subs := m.messageBroker.SubM[subscriber.Channel]
	subs = append(subs, subscriber)
	m.messageBroker.SubM[subscriber.Channel] = subs
}

type TestServer struct {
	srv *httptest.Server
	t   *testing.T
	i   controllerInterface.IController
	wg  *sync.WaitGroup
}

func TestPublishMessage(t *testing.T) {

	type TempPublisher struct {
		Id      interface{} `validate:"required"`
		Channel string      `form:"channel" json:"channel" validate:"required"`
	}

	type TempUpdates struct {
		Publisher TempPublisher `form:"publisher" json:"publisher" validate:"required"`
		Update    int           `form:"update" json:"update" validate:"required"`
	}

	tStruct := &TestServer{
		t:  t,
		wg: &sync.WaitGroup{},
	}

	tests := []struct {
		name             string
		requestBody      interface{}
		expectedResponse model.Response
		setupFunc        func(controllerInterface.IController, interface{})
		validateFunc     func(*httptest.ResponseRecorder)
	}{
		{
			name: "Success:: Publish Message :: With Encryption",
			requestBody: model.Updates{
				Publisher: model.Publisher{Id: "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e", Channel: "c4"},
				Update:    "Hello World",
			},
			setupFunc: func(i controllerInterface.IController, reqbody interface{}) {
				tStruct.i = i
				tStruct.wg.Add(1)
				testClient(tStruct, true)

				var x *models = i.(*models)
				var y model.Updates = reqbody.(model.Updates)
				m, ok := x.messageBroker.PubM[y.Publisher.Channel]
				if !ok {
					m = make(map[string]struct{})
					m[y.Publisher.Id] = struct{}{}
				}
				x.messageBroker.PubM[y.Publisher.Channel] = m
				t.Log(x.messageBroker.SubM)
			},
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = model.Response{
					Status:  http.StatusOK,
					Message: constants.NotifiedSub,
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
			},
			expectedResponse: model.Response{
				Status:  http.StatusOK,
				Message: constants.NotifiedSub,
				Data:    nil,
			},
		},
		{
			name: "Success:: Publish Message::Without Encryption",
			requestBody: model.Updates{
				Publisher: model.Publisher{Id: "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e", Channel: "c4"},
				Update:    "Hello World",
			},
			setupFunc: func(i controllerInterface.IController, reqbody interface{}) {
				tStruct.i = i
				tStruct.wg.Add(1)
				testClient(tStruct, false)

				var x *models = i.(*models)
				var y model.Updates = reqbody.(model.Updates)
				m, ok := x.messageBroker.PubM[y.Publisher.Channel]
				if !ok {
					m = make(map[string]struct{})
					m[y.Publisher.Id] = struct{}{}
				}
				x.messageBroker.PubM[y.Publisher.Channel] = m
				t.Log(x.messageBroker.SubM)
			},
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = model.Response{
					Status:  http.StatusOK,
					Message: constants.NotifiedSub,
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				t.Log(response)
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
			},
			expectedResponse: model.Response{
				Status:  http.StatusOK,
				Message: constants.NotifiedSub,
				Data:    nil,
			},
		},
		{
			name: "FAILURE::Publish Message::Incorrect UUID",
			requestBody: model.Updates{
				Publisher: model.Publisher{Id: "publisher1", Channel: "c4"},
				Update:    "Hello World",
			},
			setupFunc: func(i controllerInterface.IController, reqbody interface{}) {
			},
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = model.Response{
					Status:  http.StatusBadRequest,
					Message: constants.InvalidUUID,
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
			},
			expectedResponse: model.Response{
				Status:  http.StatusNotFound,
				Message: constants.PublisherNotFound,
				Data:    nil,
			},
		},
		{
			name: "FAILURE::Publish Message::No Publisher Found",
			requestBody: model.Updates{
				Publisher: model.Publisher{Id: "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e", Channel: "c4"},
				Update:    "Hello World",
			},
			setupFunc: func(i controllerInterface.IController, reqbody interface{}) {
			},
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = model.Response{
					Status:  http.StatusNotFound,
					Message: constants.PublisherNotFound,
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
			},
			expectedResponse: model.Response{
				Status:  http.StatusNotFound,
				Message: constants.PublisherNotFound,
				Data:    nil,
			},
		},
		{
			name:        "FAILURE::Publish Message::Incorrect input details",
			requestBody: TempUpdates{Publisher: TempPublisher{Id: "", Channel: ""}, Update: 1},
			setupFunc: func(i controllerInterface.IController, reqbody interface{}) {
			},
			validateFunc: func(w *httptest.ResponseRecorder) {
				var tempstruct = model.Response{
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
				var response model.Response
				err := json.Unmarshal(responseBody, &response)
				if err != nil {
					t.Error(error.Error())
				}
				t.Log(response)
				if !reflect.DeepEqual(response, tempstruct) {
					t.Errorf("Want: %v, Got: %v", tempstruct, response)
				}
			},
			expectedResponse: model.Response{
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
			tt.setupFunc(i, tt.requestBody)
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

	tests := []struct {
		name             string
		requestBody      interface{}
		expectedResponse model.Response
		validateFunc     func(*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "Success:: NoRouteFound",

			validateFunc: func(w *httptest.ResponseRecorder, r *http.Request) {
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Want: Content Type as %v, Got: Content Type as %v", "application/json", contentType)
				}
				expectedResponse := model.Response{
					Status:  http.StatusNotFound,
					Message: constants.NoRoute,
					Data:    nil,
				}
				if w.Code != expectedResponse.Status {
					t.Errorf("Want: %v, Got: %v", expectedResponse.Status, w.Code)
				}
				responseBody, error := ioutil.ReadAll(w.Body)
				if error != nil {
					t.Error(error.Error())
				}
				var response model.Response
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
