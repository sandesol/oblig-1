package handlers

import (
	"assignment1/consts"
	"fmt"
	"net/http"
	"strings"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	// CountriesNow url
	urlCN := consts.COUNTRIESNOWURL + "countries/population/cities"

	payload := strings.NewReader(`{
		"city": "oslo"
	}`)

	// Makes a post request and handles errors. Defers closing the body response.
	respCN, errCN := http.Post(urlCN, "application/json", payload)
	if errCN != nil {

		// -- TODO --

		fmt.Println(errCN.Error()) // debug
	}
	defer respCN.Body.Close()

	// REST Countries url
	urlRC := consts.RESTCOUNTRIESURL + "capital/oslo"

	// Makes a get request and handles errors, similarly to CountriesNow
	respRC, errRC := http.Get(urlRC)
	if errRC != nil {

		// TODO: handle error

		fmt.Println(errRC.Error())
	}
	defer respRC.Body.Close()

	// CountriesNow status code   and   REST Countries status code
	fmt.Fprint(w, "CountriesNow API status:   ", respCN.Status, "\n")
	fmt.Fprint(w, "REST Countries API status: ", respRC.Status)

}
