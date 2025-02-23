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

	urlNOW := consts.COUNTRIESNOWURL + "countries/population/cities"

	payload := strings.NewReader(`{
		"city": "oslo"
	}`)

	// Makes a post request and handles errors. Defers closing the body response.
	respNOW, errNOW := http.Post(urlNOW, "application/json", payload)
	if errNOW != nil {
		fmt.Println("Error in POST request: ", errNOW.Error()) // debug
		return
	}
	defer respNOW.Body.Close()

	urlREST := consts.RESTCOUNTRIESURL + "capital/oslo"

	// Makes a get request and handles errors, similarly to CountriesNow
	respREST, errRC := http.Get(urlREST)
	if errRC != nil {
		fmt.Println("Error in GET request: ", errRC.Error()) // debug
		return
	}
	defer respREST.Body.Close()

	status := &Status{}
	status.CountriesNowStatus = respNOW.Status
	status.RESTCountriesStatus = respREST.Status
	status.Version = "v1"
	status.Uptime = time.Now().Unix() - UptimeStart

	// Pretty-prints json from struct
	jsonStatus, errjson := json.MarshalIndent(status, "", "    ")
	if errjson != nil {
		fmt.Println("Error: ", errjson.Error())
	}
	fmt.Fprint(w, string(jsonStatus))
}
