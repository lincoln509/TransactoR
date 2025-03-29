package logging

import (
	"fmt"
	"log"
	"net/http"
)

type Logger interface {
	Info(req *http.Request, format string, args ...interface{})
	Error(req *http.Request, format string, args ...interface{})
	Warn(req *http.Request, format string, args ...interface{})
}

type DefaultLogger struct{}

func (dl *DefaultLogger) Info(req *http.Request, format string, args ...interface{}) {
	log.Printf("[INFO] %s %s - %s", req.Method, req.URL.Path, fmt.Sprintf(format, args...))
}

func (dl *DefaultLogger) Warn(req *http.Request, format string, args ...interface{}) {
	log.Printf("[WARNING] %s %s - %s", req.Method, req.URL.Path, fmt.Sprintf(format, args...))
}

func (dl *DefaultLogger) Error(req *http.Request, format string, args ...interface{}) {
	log.Printf("[ERROR] %s %s - %s", req.Method, req.URL.Path, fmt.Sprintf(format, args...))
}
