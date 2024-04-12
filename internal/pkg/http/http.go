package http

import (
	"encoding/json"
	"fmt"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ReadBody(r *http.Request, log *zap.Logger) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(constants.FailedReadRequestBody, zap.Error(err))
		return nil, pErrors.ErrReadBody
	}

	err = r.Body.Close()
	if err != nil {
		log.Error(constants.FailedCloseRequestBody, zap.Error(err))
	}

	return body, nil
}

type JSONError struct {
	Error string `json:"error"`
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	errCause := errors.Cause(err)
	httpCode, exists := pErrors.ErrorToHTTPCode(errCause)
	if !exists {
		errCause = pErrors.ErrUnknown
	}

	if pErrors.IsJSONError(errCause) {
		jsonError := JSONError{
			Error: errCause.Error(),
		}
		SendJSON(w, r, httpCode, jsonError)
	} else {
		w.WriteHeader(httpCode)
	}
}

func SendJSON(w http.ResponseWriter, r *http.Request, status int, dataStruct any) {
	dataJSON, err := json.Marshal(dataStruct)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to marshal : %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(dataJSON)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}
