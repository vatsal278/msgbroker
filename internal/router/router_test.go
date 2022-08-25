package router

import (
	"bytes"
	"encoding/json"
	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

type tempStruct struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func TestRegisterPublisher(t *testing.T) {

	tests := []struct {
		name             string
		requestBody      interface{}
		ErrorCase        bool
		expectedResponse tempStruct
	}{
		{
			name: "Success:: Router Test",
			requestBody: model.Publisher{
				Id:      "publisher1",
				Channel: "c4",
			},
			expectedResponse: tempStruct{
				Status:  http.StatusCreated,
				Message: constants.PublisherRegistration,
				Data:    nil,
			},
		},
	}
	router := Router()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			t.Log(w.Code)
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))

			router.ServeHTTP(w, r)
		})
	}
}
