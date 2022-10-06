package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func helloWorld() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		Fool bool `json:"fool"`
	}

	return WrapHandler(func(req *request) (*response, error) {
		if strings.EqualFold(req.Name, "joe") {
			return &response{
				Fool: false,
			}, nil
		}

		if strings.EqualFold(req.Name, "max") {
			return nil, APIError{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("%s: invalid name", req.Name),
			}
		}

		return &response{
			Fool: true,
		}, nil
	})
}

func main() {
	r := chi.NewRouter()

	r.Get("/", helloWorld())

	_ = http.ListenAndServe(":3030", r)
}
