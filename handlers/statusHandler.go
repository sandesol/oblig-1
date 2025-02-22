package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Status struct {
	CountriesNowStatus  string `json:"countriesnowstatus"`
	RESTCountriesStatus string `json:"restcountriesstatus"`
	Version             string `json:"version"`
	Uptime              int64  `json:"uptime"`
}

var UptimeStart int64

func InitializeUptime() {
	UptimeStart = time.Now().Unix()
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	status := &Status{}

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

		fmt.Println(errRC.Error()) // debug
	}
	defer respRC.Body.Close()

	status.CountriesNowStatus = respCN.Status
	status.RESTCountriesStatus = respRC.Status
	status.Version = "v1"
	status.Uptime = time.Now().Unix() - UptimeStart

	// Pretty-prints json from struct
	jsonStatus, errjson := json.MarshalIndent(status, "", "    ")
	if errjson != nil {
		fmt.Printf("Error: %s", errjson.Error())
	}
	fmt.Fprint(w, string(jsonStatus))
}
