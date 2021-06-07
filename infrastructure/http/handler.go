package http

import (
	"log"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request) error

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}
	// Error handling
	log.Println(err)

	clientError, ok := err.(ClientError)
	if !ok {
		// If the error is not ClientError, assume that it is ServerError.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := clientError.ResponseBody()
	if err != nil {
		log.Printf("An error ocurred: %+v", err)
		w.WriteHeader(500)
		return
	}

	status, headers := clientError.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}

type ClientError interface {
	Error() string
	// ResponseBody returns response body.
	ResponseBody() ([]byte, error)
	// ResponseHeaders returns http status code and headers.
	ResponseHeaders() (int, map[string]string)
}
