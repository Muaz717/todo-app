package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Muaz717/todo-app/internal/app/http-server/handlers/auth"
	"github.com/Muaz717/todo-app/internal/app/http-server/handlers/auth/mocks"
	resp "github.com/Muaz717/todo-app/internal/lib/api/response"
	"github.com/Muaz717/todo-app/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterNewUser(t *testing.T) {
	tests := []struct {
		name       string
		req        auth.Request
		statusCode int
		respError  string
		mockError  error
	}{
		{
			name: "Success",
			req: auth.Request{
				Email:    "test@mail.ru",
				Password: "test_password",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Empty email",
			req: auth.Request{
				Password: "test_password",
			},
			statusCode: http.StatusBadRequest,
			respError:  "field Email is a required field",
		},
		{
			name: "Empty password",
			req: auth.Request{
				Email: "test@mail.ru",
			},
			statusCode: http.StatusBadRequest,
			respError:  "field Password is a required field",
		},
		{
			name: "Invalid Email",
			req: auth.Request{
				Email:    "invalid email",
				Password: "test_password",
			},
			statusCode: http.StatusBadRequest,
			respError:  "field Email is not a valid Email",
		},
		{
			name: "RegisterNewUser error",
			req: auth.Request{
				Email:    "test@mail.ru",
				Password: "test_password",
			},
			statusCode: http.StatusInternalServerError,
			respError:  "failed to register new user",
			mockError:  errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			log := slogdiscard.NewDiscardLogger()

			authMock := mocks.NewAuth(t)

			if tt.respError == "" || tt.mockError != nil {
				authMock.
					On("RegisterNewUser", ctx, tt.req.Email, tt.req.Password).
					Return(int64(1), tt.mockError)
			}

			authHandler := auth.New(ctx, log, authMock)
			handler := authHandler.RegisterNewUser

			var input bytes.Buffer
			err := json.NewEncoder(&input).Encode(tt.req)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/sign-up/", &input)

			rr := httptest.NewRecorder()
			handler(rr, req)

			body := rr.Body.String()

			var resp resp.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.statusCode, rr.Code)

			require.Equal(t, tt.respError, resp.Error)
		})

	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		req        auth.Request
		statusCode int
		respError  string
		mockError  error
	}{
		{
			name: "Success",
			req: auth.Request{
				Email:    "test@mail.ru",
				Password: "test_password",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Empty email",
			req: auth.Request{
				Password: "test_password",
			},
			statusCode: http.StatusBadRequest,
			respError:  "field Email is a required field",
		},
		{
			name: "Empty password",
			req: auth.Request{
				Email: "test@mail.ru",
			},
			statusCode: http.StatusBadRequest,
			respError:  "field Password is a required field",
		},
		{
			name: "Invalid Email",
			req: auth.Request{
				Email:    "invalid email",
				Password: "test_password",
			},
			statusCode: http.StatusBadRequest,
			respError:  "field Email is not a valid Email",
		},
		{
			name: "Login error",
			req: auth.Request{
				Email:    "test@mail.ru",
				Password: "test_password",
			},
			statusCode: http.StatusInternalServerError,
			respError:  "invalid email or password",
			mockError:  errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			log := slogdiscard.NewDiscardLogger()

			authMock := mocks.NewAuth(t)

			if tt.respError == "" || tt.mockError != nil {
				authMock.
					On("Login", ctx, tt.req.Email, tt.req.Password).
					Return("", tt.mockError)
			}

			authHandler := auth.New(ctx, log, authMock)
			handler := authHandler.Login

			var input bytes.Buffer
			err := json.NewEncoder(&input).Encode(tt.req)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/sign-up/", &input)

			rr := httptest.NewRecorder()
			handler(rr, req)

			body := rr.Body.String()

			var resp resp.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			assert.Equal(t, rr.Code, tt.statusCode)

			assert.Equal(t, tt.respError, resp.Error)
		})

	}
}
