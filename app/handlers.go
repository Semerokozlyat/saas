package main

import (
	"net/http"
	"fmt"
	"log"
)

func (s *Service) HandlerStatus(rw http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(rw, "Server status is OK")
	if err != nil {
		log.Fatalf("failed to write a response in HandlerStatus: %v", err)
	}
	rw.WriteHeader(http.StatusOK)
}

func (s *Service) HandlerMakeScreenshot(rw http.ResponseWriter, r *http.Request) {
	s.message <- "screenshot"

	rw.WriteHeader(http.StatusAccepted)
}