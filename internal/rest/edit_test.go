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

func TestService_editUserInfo(t *testing.T) {
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
			name: "Common profile info successfully updated (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfileCommonV1{Name: "J. Doe"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.userProfileCommon.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UpdateCommonProfileInfoV1", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusOK,
		},
		{
			name: "Common profile info cannot be updated - internal error (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfileCommonV1{Name: "J. Doe"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.userProfileCommon.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UpdateCommonProfileInfoV1", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("test error"))
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Password successfully updated (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfilePasswordV1{Current: "My1stOlGoodPassword", New: "My2ndPrettyNewPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.userProfilePassword.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("ValidateCredentials", mock.Anything, mock.Anything, "My1stOlGoodPassword").Return(true, nil)
					s.On("ChangePassword", mock.Anything, mock.Anything, "My2ndPrettyNewPassword").Return(nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusOK,
		},
		{
			name: "Password cannot be updated - incorrect current (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfilePasswordV1{Current: "My1stOlGoodPassword", New: "My2ndPrettyNewPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.userProfilePassword.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("ValidateCredentials", mock.Anything, mock.Anything, "My1stOlGoodPassword").Return(false, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Password cannot be updated - incorrect new (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfilePasswordV1{Current: "My1stOlGoodPassword", New: "123"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.userProfilePassword.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Password cannot be updated - internal error (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfilePasswordV1{Current: "My1stOlGoodPassword", New: "My2ndPrettyNewPassword"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.userProfilePassword.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("ValidateCredentials", mock.Anything, mock.Anything, "My1stOlGoodPassword").Return(true, nil)
					s.On("ChangePassword", mock.Anything, mock.Anything, "My2ndPrettyNewPassword").Return(errors.New("test error"))
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Unsupported profile data (415)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					info := models.UserProfileCommonV1{Name: "J. Doe"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(info)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PATCH", "/", &buf)
					r.Header.Set("Content-Type", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.editUserInfo(tt.args.w, tt.args.r)

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
