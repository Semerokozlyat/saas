package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

func (s *Service) HandlerStatus(rw http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(rw, "Server status is OK\n")
	if err != nil {
		log.Printf("failed to write a response in HandlerStatus: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.WriteHeader(http.StatusOK)
}

func (s *Service) HandlerMakeScreenshot(rw http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	rawURL, ok := requestParams["url"]
	if !ok {
		_, _ = fmt.Fprint(rw, "Mandatory GET request parameter is absent: url\n")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	resultedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		_, _ = fmt.Fprint(rw, "URL provided is not valid\n")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	message := Message{websiteURL: resultedURL.String()}
	s.message <- message

	rw.WriteHeader(http.StatusAccepted)
}
