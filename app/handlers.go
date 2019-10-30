package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func (s *Service) HandlerStatus(rw http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(rw, "Server status is OK\n")
	if err != nil {
		log.Printf("failed to write a response in HandlerStatus: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.WriteHeader(http.StatusOK)
}

func (s *Service) HandlerRequestScreenshot(rw http.ResponseWriter, r *http.Request) {
	var bytesData []byte
	fileName, ok := r.URL.Query()["file_name"]
	if !ok {
		_, _ = fmt.Fprintf(rw, "Mandatory GET request parameter is absent: file_name\n")
		rw.WriteHeader(http.StatusBadRequest)
	}
	message := Message{screenFileName: fileName[0]}
	s.get <- message

	bytesData = DataCache[fileName[0]]

	if len(bytesData) == 0 {
		log.Println("Did not find data in cache, requesting from storage")
		bytesData, _ = s.storage.Get(fileName[0])
	}
	f, err := os.Create(fmt.Sprintf("/screens/%s.png", fileName))
	if err != nil {
		_, _ = fmt.Fprintf(rw, "Failed to create file on demand: %v", err)
	}
	bytesWritten, err := f.Write(bytesData)
	rw.Write([]byte(fmt.Sprintf("file saved successfully, its name is %s, size is %d", fileName, bytesWritten)))
	rw.WriteHeader(http.StatusOK)
}

func (s *Service) HandlerMakeScreenshot(rw http.ResponseWriter, r *http.Request) {

	rawURL, ok := r.URL.Query()["url"]
	if !ok {
		_, _ = fmt.Fprint(rw, "Mandatory GET request parameter is absent: url\n")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	resultedURL, err := url.ParseRequestURI(rawURL[0])
	if err != nil {
		_, _ = fmt.Fprint(rw, "URL provided is not valid\n")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	message := Message{websiteURL: resultedURL.String()}
	s.put <- message

	rw.WriteHeader(http.StatusAccepted)
}
