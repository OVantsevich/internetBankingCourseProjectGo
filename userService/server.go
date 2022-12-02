package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/server"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	defer repository.Close()
	if err := domain.InitConfig(); err != nil {
		log.Fatalf("environment error: %v", err)
	}
	t := server.MakeHTTPHandler()
	log.Fatal(http.ListenAndServe(":12345", t))
	repository.Close()
}
