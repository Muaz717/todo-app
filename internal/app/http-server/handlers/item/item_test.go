package item_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Muaz717/todo-app/internal/app/http-server/handlers/item"
	"github.com/Muaz717/todo-app/internal/app/http-server/handlers/item/mocks"
	"github.com/Muaz717/todo-app/internal/app/http-server/middleware/identification"
	"github.com/Muaz717/todo-app/internal/domain/models"
	resp "github.com/Muaz717/todo-app/internal/lib/api/response"
	"github.com/Muaz717/todo-app/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: .....
func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		statusCode  int
		userId      int64
		respError   string
		mockError   error
	}{
		{
			name:        "Success",
			title:       "test_title",
			description: "test_description",
			statusCode:  http.StatusOK,
			userId:      1,
		},
		{
			name:        "Empty title",
			description: "test_description",
			statusCode:  http.StatusBadRequest,
			userId:      1,
			respError:   "field Title is a required field",
		},
		{
			name:       "Empty description",
			title:      "test_title",
			statusCode: http.StatusBadRequest,
			userId:     1,
			respError:  "field Description is a required field",
		},
		{
			name:        "Create error",
			title:       "test_title",
			description: "test_description",
			userId:      1,
			statusCode:  http.StatusInternalServerError,
			respError:   "failed to create item",
			mockError:   errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			log := slogdiscard.NewDiscardLogger()

			itemHandlerMock := mocks.NewItem(t)

			if tt.respError == "" || tt.mockError != nil {
				itemHandlerMock.
					On("Create", ctx, mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(int64(1), tt.mockError)
			}

			itemHandler := item.New(ctx, log, itemHandlerMock)
			handler := itemHandler.Create

			reqBody := item.Request{
				Title:       tt.title,
				Description: tt.description,
			}

			var input bytes.Buffer
			err := json.NewEncoder(&input).Encode(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/items/", &input)

			uidStr := identification.Uid("user_id")

			withValue := context.WithValue(req.Context(), uidStr, tt.userId)
			req = req.WithContext(withValue)

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
func TestAllItemsHandler(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		userId     int64
		respError  string
		mockError  error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			userId:     1,
		},
		{
			name:       "AllItems error",
			userId:     1,
			statusCode: http.StatusInternalServerError,
			respError:  "failed to get items",
			mockError:  errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			log := slogdiscard.NewDiscardLogger()

			itemHandlerMock := mocks.NewItem(t)

			if tt.respError == "" || tt.mockError != nil {
				itemHandlerMock.
					On("AllItems", ctx, mock.AnythingOfType("int64")).
					Return([]models.Item{}, tt.mockError)
			}

			itemHandler := item.New(ctx, log, itemHandlerMock)
			handler := itemHandler.AllItems

			req := httptest.NewRequest(http.MethodGet, "/api/items/", nil)

			uidStr := identification.Uid("user_id")

			withValue := context.WithValue(req.Context(), uidStr, tt.userId)
			req = req.WithContext(withValue)

			rr := httptest.NewRecorder()
			handler(rr, req)

			require.Equal(t, tt.statusCode, rr.Code)

			if rr.Code != http.StatusOK {
				body := rr.Body.String()

				var resp resp.Response
				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				require.Equal(t, tt.respError, resp.Error)
			}

		})
	}
}
