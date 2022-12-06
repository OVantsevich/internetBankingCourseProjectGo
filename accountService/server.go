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
	err := domain.InitConfig()
	if err != nil {
		log.Fatalf("error with init config: %v", err)
	}
	eventStreaming.JetStreamInit()
	eventStreaming.HandleUser()

	h := server.MakeHTTPHandler()
	log.Fatal(http.ListenAndServe(":12344", h))
	repository.Close()
}
