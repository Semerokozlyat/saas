package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"saas/pkg/service"
)

var Version string

func main() {

	router := mux.NewRouter()
	s := service.NewService()
	fs := http.FileServer(http.Dir("/screens"))

	router.HandleFunc("/status", s.HandlerStatus).Methods(http.MethodGet)
	router.HandleFunc("/make_screenshot", s.HandlerMakeScreenshot).Methods(http.MethodGet)
	router.HandleFunc("/get_screenshot", s.HandlerRequestScreenshot).Methods(http.MethodGet)
	router.PathPrefix("/screens/").Handler(http.StripPrefix("/screens/", fs))

	go s.StartProcessing()

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

//docker container run -it --rm -v $(pwd):/usr/src/app zenika/alpine-chrome --no-sandbox --screenshot --disable-gpu --hide-scrollbars https://www.drom.ru
