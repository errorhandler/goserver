package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Err        error
}

func (e APIError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"status": e.StatusCode,
		"error":  e.Err.Error(),
	})
}

func (e APIError) Error() string {
	return e.Err.Error()
}

type Handler[Req any, Res any] func(*Req) (*Res, error)

type ErrorHandler func(writer http.ResponseWriter, request *http.Request) error

func WrapErrorHandler(handler ErrorHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			var apiError APIError

			if !errors.As(err, &apiError) {
				apiError.Err = err
				apiError.StatusCode = http.StatusInternalServerError
			}

			writer.WriteHeader(apiError.StatusCode)
			_ = json.NewEncoder(writer).Encode(apiError)

			return
		}
	}
}

func WrapHandler[Req any, Res any](handler Handler[Req, Res]) http.HandlerFunc {
	return WrapErrorHandler(func(writer http.ResponseWriter, request *http.Request) error {
		var req Req
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			return APIError{
				Err:        fmt.Errorf("failed to parse body: %w", err),
				StatusCode: http.StatusBadRequest,
			}
		}

		res, err := handler(&req)

		if err != nil {
			return err
		}

		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(writer).Encode(res); err != nil {
			return fmt.Errorf("failed to encode response: %w", err)
		}

		return nil
	})
}
