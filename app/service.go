package main

import (
	"log"
	"saas/pkg/chrome"
	"saas/pkg/storage"
	"sync"
)

type Service struct {
	message chan string
	stop    chan bool
	wg      *sync.WaitGroup

	storage storage.Storage
}

func NewService() *Service {
	service := &Service{
		message: nil,
		stop:    nil,
		wg:      nil,
	}
	service.message = make(chan string, 100)
	service.stop = make(chan bool, 1)
	service.wg = new(sync.WaitGroup)

	service.storage = storage.CreateRedisStorage()

	return service
}

func (s *Service) StartProcessing() {
	log.Println("Service started")

	//s.wg.Add(1)
	//s.wg.Wait()

	defer close(s.message)
	defer close(s.stop)

	for {
		select {
		case <-s.stop:
			log.Println("Stop processing!")
			return
		case mess, ok := <-s.message:
			if !ok {
				return
			}
			log.Printf("Get message to process: %v", mess)
			s.wg.Add(1)
			result, _ := chrome.MakeScreenshot()
			log.Println(result)
			s.wg.Done()
			s.wg.Wait()
		}
	}
}

func (s *Service) StopProcessing() {
	s.stop <- true
	s.wg.Done()
	log.Println("Service stopped")
}
