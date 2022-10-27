package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type helloWorldRequest struct {
	Name string `json:"name"`
}
type helloWorldResponse struct {
	Fool bool `json:"fool"`
}

func helloWorld() http.HandlerFunc {
	return WrapHandler(func(req *helloWorldRequest) (*helloWorldResponse, error) {
		if strings.EqualFold(req.Name, "joe") {
			return &helloWorldResponse{
				Fool: false,
			}, nil
		}

		if strings.EqualFold(req.Name, "max") {
			return nil, APIError{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("%s: invalid name", req.Name),
			}
		}

		return &helloWorldResponse{
			Fool: true,
		}, nil
	})
}

func main() {
	r := chi.NewRouter()

	r.Get("/", helloWorld())

	_ = http.ListenAndServe(":3030", r)
}
