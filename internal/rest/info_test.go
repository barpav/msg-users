package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-users/internal/rest/mocks"
	"github.com/barpav/msg-users/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_getUserInfo(t *testing.T) {
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
		wantBody    io.Reader
		wantStatus  int
	}{
		{
			name: "User info successfully received (200) - without id",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					r.Header.Set("Accept", "application/vnd.userInfo.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UserInfoV1", mock.Anything, mock.Anything).Return(
						&models.UserInfoV1{Name: "Jane Doe"},
						nil,
					)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/vnd.userInfo.v1+json",
			},
			wantBody: func() io.Reader {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(&models.UserInfoV1{Name: "Jane Doe"})
				if err != nil {
					log.Fatal(err)
				}
				return &buf
			}(),
			wantStatus: http.StatusOK,
		},
		{
			name: "User info successfully received (200) - with id",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?id=jane", nil)
					r.Header.Set("Accept", "application/vnd.userInfo.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UserInfoV1", mock.Anything, "jane").Return(
						&models.UserInfoV1{Name: "Jane Doe"},
						nil,
					)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/vnd.userInfo.v1+json",
			},
			wantBody: func() io.Reader {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(&models.UserInfoV1{Name: "Jane Doe"})
				if err != nil {
					log.Fatal(err)
				}
				return &buf
			}(),
			wantStatus: http.StatusOK,
		},
		{
			name: "User not found (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?id=jane", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UserInfoV1", mock.Anything, "jane").Return(
						nil,
						&ErrUserNotFoundTest{},
					)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Requested media type is not supported (406)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?id=jane", nil)
					r.Header.Set("Accept", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotAcceptable,
		},
		{
			name: "User deleted (410)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?id=jane", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UserInfoV1", mock.Anything, "jane").Return(
						nil,
						&ErrUserDeletedTest{},
					)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusGone,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/?id=jane", nil)
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UserInfoV1", mock.Anything, "jane").Return(
						nil,
						errors.New("test error"),
					)
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
			s.getUserInfo(tt.args.w, tt.args.r)

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

			if tt.wantBody != nil {
				require.Equal(t, tt.wantBody, tt.args.w.Body)
			}

			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}

type ErrUserNotFoundTest struct{}

func (e *ErrUserNotFoundTest) Error() string {
	return "user not found (test)"
}

func (e *ErrUserNotFoundTest) ImplementsUserNotFoundError() {
}

type ErrUserDeletedTest struct{}

func (e *ErrUserDeletedTest) Error() string {
	return "user deleted (test)"
}

func (e *ErrUserDeletedTest) ImplementsUserDeletedError() {
}
