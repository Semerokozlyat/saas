package service

import (
	"log"
	"saas/pkg/chrome"
	"saas/pkg/storage"
	"sync"
	"time"
)

var (
	DataCache = make(map[string][]byte)
)

type Service struct {
	get     chan Message
	put     chan Message
	stop    chan bool
	wg      *sync.WaitGroup

	storage storage.Storage
}

type Message struct {
	ScreenFileName string
	WebsiteURL     string
}

func NewService() *Service {
	service := &Service{
		put: nil,
		stop:    nil,
		wg:      nil,
	}
	service.get = make(chan Message, 100)
	service.put = make(chan Message, 100)
	service.stop = make(chan bool, 1)
	service.wg = new(sync.WaitGroup)

	service.storage = storage.CreateRedisStorage()

	return service
}

func (s *Service) StartProcessing() {
	log.Println("Service started")

	//s.wg.Add(1)
	//s.wg.Wait()

	defer close(s.put)
	defer close(s.stop)

	ticker := time.NewTicker(300 * time.Second)

	for {
		select {
		case <-s.stop:
			log.Println("Stop processing!")
			return
		case <- ticker.C:
			DataCache = make(map[string][]byte)
		case mess, ok := <-s.get:
			if !ok {
				return
			}
			log.Printf("Received request for file: %v", mess)
			s.wg.Add(1)
			fileByteData, errGet := s.storage.Get(mess.ScreenFileName)
			if errGet != nil {
				log.Printf("failed to get file data: %v", errGet)
			}
			DataCache[mess.ScreenFileName] = fileByteData
			s.wg.Done()
			s.wg.Wait()
		case mess, ok := <-s.put:
			if !ok {
				return
			}
			log.Printf("Received message to process: %v", mess)
			s.wg.Add(1)
			fileByteData, err := chrome.MakeScreenshot(mess)
			if err != nil {
				log.Printf("failed to make screenshot: %v", err)
			}
			errPut := s.storage.Put(mess.ScreenFileName, fileByteData)
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
