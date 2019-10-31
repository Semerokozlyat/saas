package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"saas/pkg/service"
)

var Version string

const AppName string = "Screenshot-as-a-Service"

func main() {

	var versionFlag = flag.Bool("--version", false, "show application version")

	flag.Parse()
	if *versionFlag {
		fmt.Println("Application name: ", AppName)
		fmt.Println("Application version: ", Version)
	}

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
