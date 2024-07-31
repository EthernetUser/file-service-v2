package save_test

import (
	"bytes"
	"context"
	"file-service/m/internal/handlers/save"
	"file-service/m/internal/handlers/save/mocks"
	mockLogger "file-service/m/internal/logger/mocks"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveHandler(t *testing.T) {
	log := mockLogger.NewLogger()
	db := mocks.NewDb(t)
	storage := mocks.NewStorage(t)
	uuidGen := mocks.NewUuidGenerator(t)
	error := fmt.Errorf("error")
	handler := save.New(log, db, storage, uuidGen)
	uuidGen.On("GenerateUUID").Return("123").Maybe()
	storage.On("GetStoragePath").Return("test").Maybe()
	storage.On("GetStorageType").Return("local").Maybe()

	t.Run("success", func(t *testing.T) {
		storage.On("SaveFile", mock.Anything, mock.Anything).Return(nil).Once()
		db.On("SaveFile", mock.Anything).Return(int64(1), nil).Once()

		r, w := CreateRequestAndResponse(t, []byte("test"), "file", "test")

		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body);

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"success\",\"message\":\"file saved\",\"id\":1}\n", bodyResp)
	})

	t.Run("invalid file key", func(t *testing.T) {
		r, w := CreateRequestAndResponse(t, []byte("test"), "invalid file key", "test")
		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body);

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"invalid request\"}\n", bodyResp)
	})

	t.Run("invalid file name", func(t *testing.T) {
		r, w := CreateRequestAndResponse(t, []byte("test"), "file", "")
		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body);

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"invalid request\"}\n", bodyResp)
	})

	t.Run("storage error", func(t *testing.T) {
		storage.On("SaveFile", mock.Anything, mock.Anything).Return(error).Once()
		r, w := CreateRequestAndResponse(t, []byte("test"), "file", "test")
		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body);

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"failed to save file\"}\n", bodyResp)
	})

	t.Run("db error", func(t *testing.T) {
		storage.On("SaveFile", mock.Anything, mock.Anything).Return(nil).Once()
		db.On("SaveFile", mock.Anything).Return(int64(0), error).Once()
		r, w := CreateRequestAndResponse(t, []byte("test"), "file", "test")
		handler.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		bodyResp := string(body);

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "{\"status\":\"error\",\"message\":\"failed to save file\"}\n", bodyResp)
	})
}

func CreateRequestAndResponse(t *testing.T, file []byte, fileKey string, fileName string) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)

	filePart, _ := multipartWriter.CreateFormFile(fileKey, fileName)
	filePart.Write(file)
	multipartWriter.Close()

	r := httptest.NewRequest("POST", "/", &buf)
	r.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	r = r.WithContext(
		context.WithValue(r.Context(), "requestId", "123"),
	)
	w := httptest.NewRecorder()

	return r, w
}
