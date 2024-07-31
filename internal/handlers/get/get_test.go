package get_test

import (
	"context"
	"file-service/m/internal/database"
	"file-service/m/internal/handlers/get"
	"file-service/m/internal/handlers/get/mocks"
	mockLogger "file-service/m/internal/logger/mocks"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHandler(t *testing.T) {
	log := mockLogger.NewLogger()
	db := mocks.NewDb(t)
	storage := mocks.NewStorage(t)
	errorResp := fmt.Errorf("error")
	handler := get.New(log, db, storage)

	t.Run("success", func(t *testing.T) {
		db.On("GetFile", mock.Anything, mock.Anything).Return(&database.File{Id: 1}, nil).Once()
		storage.On("GetFile", mock.Anything).Return([]byte("test"), nil).Once()

		r, w := CreateRequestAndResponse("fileID", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []byte("test"), body)
	})

	t.Run("db error", func(t *testing.T) {
		db.On("GetFile", mock.Anything, mock.Anything).Return(nil, errorResp).Once()

		r, w := CreateRequestAndResponse("fileID", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"failed to get file\"}\n", bodyResp)
	})

	t.Run("storage error", func(t *testing.T) {
		db.On("GetFile", mock.Anything, mock.Anything).Return(&database.File{Id: 1}, nil).Once()
		storage.On("GetFile", mock.Anything, mock.Anything).Return(nil, errorResp).Once()

		r, w := CreateRequestAndResponse("fileID", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"failed to get file\"}\n", bodyResp)
	})

	t.Run("invalid file id", func(t *testing.T) {
		r, w := CreateRequestAndResponse("fileID", "1sdf")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"invalid file id\"}\n", bodyResp)
	})

	t.Run("invalid file key", func(t *testing.T) {
		r, w := CreateRequestAndResponse("fileKey", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"file id is empty\"}\n", bodyResp)
	})
}

func CreateRequestAndResponse(fileKey string, fileId string) (*http.Request, *httptest.ResponseRecorder) {
	r := httptest.NewRequest("Get", fmt.Sprintf("/%s", fileId), nil)
	r = r.WithContext(
		context.WithValue(r.Context(), fileKey, fileId),
	)
	r = r.WithContext(
		context.WithValue(r.Context(), "requestId", "123"),
	)
	w := httptest.NewRecorder()

	return r, w
}
