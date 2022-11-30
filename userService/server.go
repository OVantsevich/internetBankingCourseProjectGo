package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/server"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	t := server.MakeHTTPHandler()
	log.Fatal(http.ListenAndServe(":12345", t))
	repository.Close()
}
