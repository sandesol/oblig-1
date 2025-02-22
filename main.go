package main

import (
	"assignment1/consts"
	"assignment1/handlers"
	"log"
	"net/http"
)

func main() {

	handlers.InitializeUptime()

	http.HandleFunc("/countryinfo/v1/info/{two_letter_country_code}", handlers.InfoHandler)
	http.HandleFunc("/countryinfo/v1/status", handlers.StatusHandler)

	err := http.ListenAndServe(":"+consts.PORT, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
