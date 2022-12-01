package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/eventStreaming"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/server"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	_, err := domain.InitConfig()
	if err != nil {
		log.Fatalf("error with init config: %v", err)
	}

	err = eventStreaming.JetStreamInit()
	if err != nil {
		log.Errorf("error with init stream: %v", err)
	}
	_, err = eventStreaming.JetStream.Subscribe(eventStreaming.SubjectNameUserCreated, eventStreaming.CreatingUserHandler)
	if err != nil {
		log.Printf("Subscribe for " + eventStreaming.SubjectNameUserCreated + " failed")
	}

	h := server.MakeHTTPHandler()
	log.Fatal(http.ListenAndServe(":12344", h))
	repository.Close()
}
