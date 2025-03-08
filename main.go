package main

import (
	"assignment1/consts"
	"assignment1/handlers"
	"log"
	"net/http"
)

/**
 * Main function.
 *
 * - Initializes uptime
 * - Serves paths .../population/, .../info/, and .../status
 * - Starts the service and handles errors
 */
func main() {

	handlers.InitializeUptime()

	http.HandleFunc("/countryinfo/v1/population/{two_letter_country_code}", handlers.PopulationHandler)
	http.HandleFunc("/countryinfo/v1/info/{two_letter_country_code}", handlers.InfoHandler)
	http.HandleFunc("/countryinfo/v1/status", handlers.StatusHandler)

	err := http.ListenAndServe(":"+consts.PORT, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
