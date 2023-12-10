package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
)

// unmarshalBody(body io.ReadCloser, form interface{}) error
// конвертирует данные в json
func (h *APIHandler) unmarshalBody(body io.ReadCloser, form interface{}) error {

	defer body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), form); err != nil {
		return err
	}

	return nil
}

// getTextBody(body io.ReadCloser) (string, error)
// возвращает текст из тела запроса
func (h *APIHandler) getTextBody(body io.ReadCloser) (string, error) {

	defer body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getBodyData(w http.ResponseWriter, r *http.Request, data BaseBodyDater) bool
// записывает данные из тела в переменную и проверяет её валидность
// в случае не удачи записывает статус BadRequest
// возвращает true или false
func (h *APIHandler) getBodyData(w http.ResponseWriter, r *http.Request, data BaseBodyDater) bool {

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return false
	}

	if ok := data.isValid(); !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return false
	}

	return true

}

// setHTTPError(w http.ResponseWriter, err error) присваивает response статус ответа
// вынесен для исключения дублирования в коде
func (h *APIHandler) setHTTPError(w http.ResponseWriter, err error) {
	var serviceError *lErrors.ServicesError
	if errors.As(err, &serviceError) {
		http.Error(w, serviceError.Error(), serviceError.HTTPCode)
	} else {
		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
	}
}
