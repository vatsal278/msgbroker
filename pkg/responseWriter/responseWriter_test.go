package responseWriter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name             string
		expectedResponse interface{}
	}{
		{
			name:             "SUCCESS:: ParseRequest",
			expectedResponse: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run("TestValidate", func(t *testing.T) {
			w := httptest.NewRecorder()
			//r := httptest.
			//Mock the interface and use them inside it
			err := responseWriter.ResponseWriter(w, http.StatusOK, "Successfully Registered as publisher to the channel", nil, &model.Response{})
			if err != nil {
				t.Errorf("Want: %v, Got: %v", nil, err.Error())
			}
			if w.Code != http.StatusOK {
				t.Errorf("Want: %v, Got: %v", tt.expectedResponse, w.Code)
			}
		},
		)
	}
}
