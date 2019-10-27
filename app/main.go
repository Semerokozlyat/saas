package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "saas/pkg/chrome"
	_ "time"
)

var Version string

func main() {

	router := mux.NewRouter()
	s := NewService()

	router.HandleFunc("/status", s.HandlerStatus).Methods(http.MethodGet)
	router.HandleFunc("/make_screenshot", s.HandlerMakeScreenshot).Methods(http.MethodPost)

	go s.StartProcessing()

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

//docker container run -it --rm -v $(pwd):/usr/src/app zenika/alpine-chrome --no-sandbox --screenshot --disable-gpu --hide-scrollbars https://www.drom.ru
