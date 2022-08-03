package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	tests := []struct {
		name              string
		requestBody       interface{}
		ErrorCase         bool
		expected_response temp_struct
	}{
		{
			name:        "Success:: Router Test",
			requestBody: publisher,
			expected_response: temp_struct{
				Status:  http.StatusCreated,
				Message: "Successfully Registered as publisher to the channel",
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
