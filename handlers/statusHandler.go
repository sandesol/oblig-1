package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

/**
 * A struct that contains the status of both apis, version and uptime
 */
type Status struct {
	CountriesNowStatus  string `json:"countriesnowstatus"`
	RESTCountriesStatus string `json:"restcountriesstatus"`
	Version             string `json:"version"`
	Uptime              int64  `json:"uptime"`
}

// variable for uptime
var UptimeStart int64

/**
 * Tiny function called by the main function that initializes the uptime variable
 */
func InitializeUptime() {
	UptimeStart = time.Now().Unix()
}

/**
 * Driver function that is called by main.go
 *
 * @param w http.ResponseWriter - used to print json and error messages to the user
 * @param r *http.Request       - used to get request method
 */
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	// Only allows GET methods
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" method is not allowed. Use "+http.MethodGet+" method instead.", http.StatusMethodNotAllowed)
		return
	}

	urlNOW := consts.COUNTRIESNOWURL + "countries/population/cities"

	// The contents of the post request.
	// It is intentionally small to reduce traffic to the API it is sent to
	payload := strings.NewReader(`{"city": "oslo"}`)

	// Makes a post request and handles errors. Defers closing the body response.
	respNOW, errNOW := http.Post(urlNOW, "application/json", payload)
	if errNOW != nil {
		fmt.Println("Error in POST request: ", errNOW.Error()) // debug
		return
	}
	defer respNOW.Body.Close()

	urlREST := consts.RESTCOUNTRIESURL + "alpha/gn?fields=population"

	// Makes a get request and handles errors, similarly to CountriesNow
	respREST, errRC := http.Get(urlREST)
	if errRC != nil {
		fmt.Println("Error in GET request: ", errRC.Error()) // debug
		return
	}
	defer respREST.Body.Close()

	// Sets up a struct that contains the fetched status codes, a hard coded version
	//  and the time since the service booted
	status := Status{}
	status.CountriesNowStatus = respNOW.Status
	status.RESTCountriesStatus = respREST.Status
	status.Version = "v1"
	status.Uptime = time.Now().Unix() - UptimeStart

	// Pretty-prints json from the struct we just set up
	jsonStatus, errjson := json.MarshalIndent(status, "", "    ")
	if errjson != nil {
		fmt.Println("Error: ", errjson.Error())
	}
	fmt.Fprint(w, string(jsonStatus))
}
