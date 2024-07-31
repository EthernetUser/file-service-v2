package setdelete_test

import (
	"context"
	"file-service/m/internal/handlers/setdelete/mocks"
	setdelete "file-service/m/internal/handlers/setDelete"
	mockLogger "file-service/m/internal/logger/mocks"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetDeleteHandler(t *testing.T) {
	log := mockLogger.NewLogger()
	db := mocks.NewDb(t)
	errorResp := fmt.Errorf("error")
	handler := setdelete.New(log, db)

	t.Run("success", func(t *testing.T) {
		db.On("SetFileIsDeleted", mock.Anything).Return(int64(1), nil).Once()

		r, w := CreateRequestAndResponse("fileID", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"success\",\"message\":\"file deleted\"}\n", bodyResp)
	})

	t.Run("db error", func(t *testing.T) {
		db.On("SetFileIsDeleted", mock.Anything).Return(int64(0), errorResp).Once()

		r, w := CreateRequestAndResponse("fileID", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"failed to delete file\"}\n", bodyResp)
	})

	t.Run("0 affected rows", func(t *testing.T) {
		db.On("SetFileIsDeleted", mock.Anything).Return(int64(0), nil).Once()

		r, w := CreateRequestAndResponse("fileID", "1")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"failed to delete file\"}\n", bodyResp)
	})

	t.Run("invalid fileId", func(t *testing.T) {
		r, w := CreateRequestAndResponse("fileID", "1asdf")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"invalid file id\"}\n", bodyResp)
	})

	t.Run("invalid fileKey", func(t *testing.T) {
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
	r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", fileId), nil)
	r = r.WithContext(
		context.WithValue(r.Context(), fileKey, fileId),
	)
	r = r.WithContext(
		context.WithValue(r.Context(), "requestId", "123"),
	)
	w := httptest.NewRecorder()

	return r, w
}
