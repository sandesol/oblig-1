package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/**
 * The driver function that is called by main.go.
 *
 * @param w http.ResponseWriter - used to print json and error messages to the user
 * @param r *http.Request       - used to get request methods, and the limit from the url
 */
func PopulationHandler(w http.ResponseWriter, r *http.Request) {

	// Only allows GET methods
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" method is not allowed. Use "+http.MethodGet+" method instead.", http.StatusMethodNotAllowed)
		return
	}

	iso := r.PathValue("two_letter_country_code") // get iso code
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 must be a 2 letter code. (Error code 2000)", http.StatusBadRequest) // 400
		return
	}
	iso3, err := GetCountry(w, iso) // get country name
	if err != nil {
		return
		// Errors are handled in GetCountry().
		// The error message is just a dummy "".
		// This is done so we can check for errors and exit PopulationHandler() properly.
		// But in this case we don't really care about the message as long as we have an error.
		// Just checking for an empty string ("") feels wrong.
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" { // no args, use whole range

		errPopNoArgs := FetchPopulation(w, iso3, "", "")
		if errPopNoArgs != nil {
			fmt.Fprintln(w, errPopNoArgs.Error()) // TODO
			return
		}

	} else {
		timeframe := strings.Split(limit, "-")
		if len(timeframe) != 2 {
			http.Error(w, "Expected 2 arguments, got "+fmt.Sprint(len(timeframe))+". (Error code 2002)", http.StatusBadRequest) // :)
			return                                                                                                              // need 2 args
		}
		if timeframe[0] == "" || timeframe[1] == "" {
			http.Error(w, "One or more arguments are empty. (Error code 2003)", http.StatusBadRequest) // :)
			return                                                                                     // one or more empty args
		}

		errPopWithArgs := FetchPopulation(w, iso3, timeframe[0], timeframe[1])
		if errPopWithArgs != nil {
			return
		}
	}

}

/**
 * Gets the iso-3 code based on the valid iso-2 code that was passed as a parameter.
 *
 * @param w   http.ResponseWriter - used to write errors
 * @param iso string              - syntactically correct iso-2 code (not necessarily in use though)
 *
 * @returns - (string, error) tuple. On error, string is empty and error is an empty string. On OK,
 *              string is an iso-3 code and error is nil
 */
func GetCountry(w http.ResponseWriter, iso string) (string, error) {
	var country struct {
		Iso3 string `json:"cca3"`
	}

	resp, errGet := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + "?fields=cca3")
	if errGet != nil {
		log.Println("(FetchCountry) Error in http.Get: ", errGet.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return "", errors.New("")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Error: iso-2 code is not in use. (Error code 3000)", http.StatusNotFound)
		return "", errors.New("")
	}

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Println("(FetchCountry) Error in io.ReadAll: ", errReadAll.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return "", errors.New("")
	}

	errJson := json.Unmarshal(body, &country)
	if errJson != nil {
		log.Println("(FetchCountry) Error parsing json with json.Unmarshal: ", errJson.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return "", errors.New("")
	}

	if country.Iso3 == "" {
		http.Error(w, "Could not retrieve an iso3 code from iso2 code \""+iso+"\". (Error code 2001)", http.StatusNotFound) // 404
		return "", errors.New("")
	}

	return country.Iso3, nil
}

/**
 * Fetches the population data into a temporary struct from 'min' to 'max', and prints it out.
 * Returns an empty error if something went wrong, because we are only interested in checking
 *  for errors and not in the actual error message itself as the message is handled by http.Error().
 *
 * @param w    http.Responsewriter - for sending http errors
 * @param iso3 string              - iso-3 code of country of interest
 * @param min  string              - first year we want to get
 * @param max  string              - last year we want to get
 *
 * @returns - error with an empty message if there are any, nil otherwise
 */
func FetchPopulation(w http.ResponseWriter, iso3, min, max string) error {
	var start, end int

	if min != "" {
		s, errConvStart := strconv.Atoi(min)
		if errConvStart != nil {
			http.Error(w, "Start year must be a number. (Error code 2004.1)", http.StatusBadRequest) //  400
			return errors.New("")
		}
		start = s
	} else {
		start = 0
	}
	if max != "" {
		e, errConvEnd := strconv.Atoi(max)
		if errConvEnd != nil {
			http.Error(w, "End year must be a number. (Error code 2004.2)", http.StatusBadRequest) // 400
			return errors.New("")
		}
		end = e
	} else {
		end = time.Now().Year()
	}
	if start > end {
		http.Error(w, "Start year is greater than end year. (Error code 2005)", http.StatusBadRequest) // 400
		return errors.New("")
	}

	var wrapper struct {
		Mean int `json:"mean"`
		Data struct {
			PopulationCounts []struct {
				Year  int `json:"year"`
				Value int `json:"value"`
			} `json:"populationcounts"`
		} `json:"data"`
	}

	payload := strings.NewReader(`{"iso3": "` + iso3 + `"}`)

	// Makes a post request and handles errors. Defers closing the body response.
	resp, errPost := http.Post(consts.COUNTRIESNOWURL+"countries/population", "application/json", payload)

	if errPost != nil {
		log.Println("(FetchPopulation) Error in post request: ", errPost.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return errors.New("")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Error: iso-2 code is not in use. (Error code 3000)", http.StatusNotFound)
	}

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Println("(FetchPopulation) Error in io.ReadAll: ", errReadAll.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return errors.New("")
	}

	errJson := json.Unmarshal(body, &wrapper)
	if errJson != nil {
		log.Println("(FetchCities) There was an error parsing json: ", errJson.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return errors.New("")
	}

	var i, j = 0, 0
	// finds first instance that matches
	for ; i < len(wrapper.Data.PopulationCounts); i++ {
		if start <= wrapper.Data.PopulationCounts[i].Year {
			break
		}
	}
	for j = len(wrapper.Data.PopulationCounts) - 1; 0 <= j; j-- {
		if end >= wrapper.Data.PopulationCounts[j].Year {
			break
		}
	}

	wrapper.Data.PopulationCounts = wrapper.Data.PopulationCounts[i : j+1]

	// calculates sum of all years. 'val' is a struct, which is why we do val.Value
	var sum = 0
	for _, val := range wrapper.Data.PopulationCounts {
		sum += val.Value
	}

	// handle division by 0 error if no cities were found
	if len(wrapper.Data.PopulationCounts) != 0 {
		wrapper.Mean = sum / len(wrapper.Data.PopulationCounts)
	} else {
		wrapper.Mean = 0
	}

	jsonStatus, errjson := json.MarshalIndent(wrapper, "", "    ")
	if errjson != nil {
		log.Println("Error: ", errjson.Error())
	}
	fmt.Fprintln(w, string(jsonStatus))

	return nil
}
