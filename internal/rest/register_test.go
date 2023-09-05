package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-users/internal/rest/mocks"
	"github.com/barpav/msg-users/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_registerNewUser(t *testing.T) {
	type testService struct {
		storage Storage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		testService testService
		args        args
		wantHeaders map[string]string
		wantStatus  int
	}{
		{
			name: "New user successfully registered (201)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "jane", Name: "Jane Doe", Password: "My1stGoodPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("CreateUser", mock.Anything, "jane", "Jane Doe", "My1stGoodPassword").Return(nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusCreated,
		},
		{
			name: "Incorrect parameters (400) - user id",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "&^@%$", Name: "Jane Doe", Password: "My1stGoodPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Incorrect parameters (400) - name",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "jane", Password: "My1stGoodPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Incorrect parameters (400) - password",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "jane", Name: "Jane Doe", Password: "123"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Incorrect parameters (400) - schema",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := struct{ Anything string }{Anything: "test"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "User id already exists (409)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "jane", Name: "Jane Doe", Password: "My1stGoodPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("CreateUser", mock.Anything, "jane", "Jane Doe", "My1stGoodPassword").Return(&ErrUserIdAlreadyExistsTest{})
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusConflict,
		},
		{
			name: "Unsupported new user data (415)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "jane", Name: "Jane Doe", Password: "My1stGoodPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					newUser := models.NewUserV1{Id: "jane", Name: "Jane Doe", Password: "My1stGoodPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(newUser)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newUser.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("CreateUser", mock.Anything, "jane", "Jane Doe", "My1stGoodPassword").Return(errors.New("test error"))
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.registerNewUser(tt.args.w, tt.args.r)

			for k, v := range tt.wantHeaders {
				require.Equal(t, v, func() string {
					h := tt.args.w.Result().Header
					if h == nil {
						return ""
					}
					v := h[k]
					if len(v) == 0 {
						return ""
					}
					return v[0]
				}())
			}

			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}

type ErrUserIdAlreadyExistsTest struct{}

func (e *ErrUserIdAlreadyExistsTest) Error() string {
	return "user id already exists"
}

func (e *ErrUserIdAlreadyExistsTest) ImplementsUserIdAlreadyExistsError() {
}
