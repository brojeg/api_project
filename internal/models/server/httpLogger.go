package model

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

type key string

const correlationIDKey key = "correlationID"

var path string

func SetServerLogPath(p string) {
	path = p
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logFile, err := createLogFile()
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	logger := log.New(logFile, "", log.LstdFlags)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := generateCorrelationID()
		ctx := context.WithValue(r.Context(), correlationIDKey, correlationID)
		// create a new buffer to read the request body
		requestBodyBuffer := new(bytes.Buffer)
		teeReader := io.TeeReader(r.Body, requestBodyBuffer)
		// replace the request body with the buffer reader
		r.Body = io.NopCloser(teeReader)
		logger.Printf("%s %s %s [CorrelationID: %s]", r.RemoteAddr, r.Method, r.URL, correlationID)
		// create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		// call the next handler in the chain, passing the response recorder instead of the original response writer
		next.ServeHTTP(recorder, r.WithContext(ctx))
		// read the response body
		responseBody := recorder.Body.Bytes()
		// get the correlation ID from the context
		correlationID = ctx.Value(correlationIDKey).(string)
		logger.Printf("%s %s %s %d [CorrelationID: %s]", r.RemoteAddr, r.Method, r.URL, recorder.Code, correlationID)
		logger.Printf("Response body: %s", responseBody)
		// copy the headers from the recorder to the original response writer
		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		// set the correlation ID header on the outgoing response
		w.Header().Set("CorrelationID", correlationID)
		// set the status code on the original response writer
		w.WriteHeader(recorder.Code)
		// write the response body to the original response writer
		w.Write(responseBody)
	})
}

func createLogFile() (*os.File, error) {
	currentDate := time.Now().Format("2006-01-02")
	mydir, _ := os.Getwd()
	fmt.Println(mydir)
	// create the log file with the current date as the file name
	logFile, err := os.OpenFile(path+"server-"+currentDate+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error creating log file: %v", err)
	}

	return logFile, nil
}

func generateCorrelationID() string {
	// generate a unique identifier using the current time
	return time.Now().Format("20060102-150405.999999")
}
