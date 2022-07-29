package responseWriter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name               string
		ExpectedStatusCode int
		testcase           int
		setupFunc          func(w *httptest.ResponseRecorder) error
	}{
		{
			name:     "SUCCESS:: validate",
			testcase: 1,
			setupFunc: func(w *httptest.ResponseRecorder) error {
				err := responseWriter.ResponseWriter(w, http.StatusOK, "Successfully Registered as publisher to the channel", nil, &model.Response{})
				return err
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:     "FAILURE:: validate",
			testcase: 2,
			setupFunc: func(w *httptest.ResponseRecorder) error {
				err := responseWriter.ResponseWriter(w, http.StatusOK, "Successfully Registered as publisher to the channel", nil, &model.Response{})
				return err
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			//r := httptest.
			_ = tt.setupFunc(w)
			if tt.testcase == 1 {
				if w.Code != tt.ExpectedStatusCode {
					t.Errorf("Want: %v, Got: %v", tt.ExpectedStatusCode, w.Code)
				}
				return
			}
			if w.Code != tt.ExpectedStatusCode {
				t.Errorf("Want: %v, Got: %v", tt.ExpectedStatusCode, w.Code)
			}
		})
	}
}
