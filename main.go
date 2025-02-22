package main

import (
	"assignment1/consts"
	"assignment1/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/status", handlers.StatusHandler) // TODO: endre path (/countryinfo/v1/status/)

	err := http.ListenAndServe(":"+consts.PORT, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
