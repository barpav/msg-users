package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-users/internal/rest/mocks"
	"github.com/barpav/msg-users/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_deleteUser(t *testing.T) {
	type testService struct {
		storage Storage
		auth    Authenticator
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
			name: "User deletion request accepted (202)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("GenerateUserDeletionCode", mock.Anything, mock.Anything).Return("test-code", nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"confirmation-code": "test-code",
			},
			wantStatus: http.StatusAccepted,
		},
		{
			name: "Failed to generate confirmation code (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/", nil)
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("GenerateUserDeletionCode", mock.Anything, mock.Anything).Return("", errors.New("test error"))
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "User successfully deleted (204)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/?confirmation-code=test-code", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("ValidateUserDeletionCode", mock.Anything, mock.Anything, "test-code").Return(true, nil)
					s.On("UserInfoV1", mock.Anything, mock.Anything).Return(&models.UserInfoV1{Name: "Jane Doe"}, nil)
					s.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)
					return s
				}(),
				auth: func() *mocks.Authenticator {
					a := mocks.NewAuthenticator(t)
					a.On("EndAllSessions", mock.Anything, mock.Anything).Return(nil)
					return a
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNoContent,
		},
		{
			name: "Invalid confirmation code (417)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/?confirmation-code=test-code", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("ValidateUserDeletionCode", mock.Anything, mock.Anything, "test-code").Return(false, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusExpectationFailed,
		},
		{
			name: "Failed to delete user (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/?confirmation-code=test-code", nil)
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("ValidateUserDeletionCode", mock.Anything, mock.Anything, "test-code").Return(true, nil)
					s.On("UserInfoV1", mock.Anything, mock.Anything).Return(&models.UserInfoV1{Name: "Jane Doe"}, nil)
					s.On("DeleteUser", mock.Anything, mock.Anything).Return(errors.New("test-error"))
					return s
				}(),
				auth: func() *mocks.Authenticator {
					a := mocks.NewAuthenticator(t)
					a.On("EndAllSessions", mock.Anything, mock.Anything).Return(nil)
					return a
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
				auth:    tt.testService.auth,
			}
			s.deleteUser(tt.args.w, tt.args.r)

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
