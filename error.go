package utils

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func RaiseError(err error, code int) {
	if err != nil {
		panic(&serviceError{Message: err.Error(), Code: code})
	}
}

func RecoverIfError(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		err := rec.(error)
		log.Printf("Detected error: %s\n", err.Error())
		sendErrorResponse(w, r, err)
	}
}

// Service Error
type serviceError struct {
	Code    int
	Message string
}

func (e *serviceError) Error() string {
	return e.Message
}

// ErrorResponse
type errorResponse struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp"`
}

func (resp *errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	resp.Timestamp = time.Now().Unix()
	return nil
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusInternalServerError
	message := err.Error()

	serr, ok := err.(*serviceError)
	if ok {
		code = serr.Code
		message = serr.Message
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	render.Render(w, r, &errorResponse{Message: message, Code: code})
}
