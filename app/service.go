package main

import (
	"log"
	"saas/pkg/chrome"
	"saas/pkg/storage"
	"sync"
)

type Service struct {
	message chan Message
	stop    chan bool
	wg      *sync.WaitGroup

	storage storage.Storage
}

type Message struct {
	websiteURL string
}

func NewService() *Service {
	service := &Service{
		message: nil,
		stop:    nil,
		wg:      nil,
	}
	service.message = make(chan Message, 100)
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
			fileByteData, err := chrome.MakeScreenshot(mess.websiteURL)
			if err != nil {
				log.Printf("failed to make screenshot: %v", err)
			}
			errPut := s.storage.Put(fileByteData)
			if errPut != nil {
				log.Printf("failed to save file data: %v", errPut)
			}
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
