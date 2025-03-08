package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

/**
 * Main struct that contains all information about everything
 */
type Country struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    []string          `json:"capital"`
	Cities     []string          `json:"data"`
}

/**
 * A struct that just contains cities.
 * This struct exists solely to make filtering operations more convenient.
 */
type JutsCities struct {
	Cities []string `json:"data"`
}

/**
 * The driver function that is called by main.go.
 *
 * @param w http.ResponseWriter - used to print json and error messages to the user
 * @param r *http.Request       - used to get request methods, and the limit from the url
 */
func InfoHandler(w http.ResponseWriter, r *http.Request) {

	// Only allows GET methods
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" method is not allowed. Use "+http.MethodGet+" method instead.", http.StatusMethodNotAllowed)
		return
	}

	country := Country{}
	iso := r.PathValue("two_letter_country_code")
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 must be a 2 letter code. (Error code 1000)", http.StatusBadRequest)
		return
	}

	// fetches countries first, so we can use country name as a parameter in the post request in FetchCities()
	errFetch := FetchCountry(w, &country, iso)
	if errFetch != nil {
		return
		// we do not care about the error message as the error has already been displayed to the user.
		// we simply return an erorr to properly exit out of the handler
	}

	query := r.URL.Query().Get("limit")
	if query == "" {

		// No limit set, defaults to 10
		errCities := FetchCities(w, &country, 10)
		if errCities != nil {
			return
		}

	} else {

		// gets limit, converts it to int
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil {
			http.Error(w, "Error: limit must be an integer. (Error code 1001)", http.StatusBadRequest)
			return
		}
		if limit < 0 {
			http.Error(w, "Error: limit must be a positive integer. (Error code 1002)", http.StatusBadRequest)
			return
		}

		// passed both tests, limit is valid and we can fetch cities
		errCities := FetchCities(w, &country, limit)
		if errCities != nil {
			return
		}

	}

	// prettyprints the complete country
	PrintCountry(w, country)
}

/**
 * Fetches everything except cities into the Country pointer 'c'.
 * For some odd reason, r.PathValue can only be read once, which means we have to pass iso as a parameter
 *
 * @param w   http.ResponseWriter - writes errors to the user if there are any
 * @param c   *Country            - country passed by reference, because we want to manipulate the original
 * @param iso string 			  - iso-2 code of the country we want to fetch
 *
 * @returns - an error if there are any, nil otherwise
 */
func FetchCountry(w http.ResponseWriter, c *Country, iso string) error {

	resp, errGet := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + consts.RESTCOUNTRIESFILTER)
	if errGet != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("(FetchCountry) Error in http.Get: ", errGet.Error()) // debug
		return errors.New("")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Error: iso-2 code is not in use. (Error code 3000)", http.StatusNotFound)
		return errors.New("")
	}

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("(FetchCountry) Error in io.ReadAll: ", errReadAll.Error()) // debug
		return errors.New("")
	}

	errJson := json.Unmarshal(body, c)
	if errJson != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("(FetchCountry) There was an error parsing json: ", errJson.Error())
		return errors.New("")
	}

	return nil
}

/**
 *
 * Fetches all cities via countriesNow into a temporary struct, filters the temporary struct based on parameters,
 * copies the contents of the filtered cities into the "main struct" and finally sorts the cities in ascending order.
 *
 * The fetched cities from CountriesNow are only partially sorted, so we need to manually sort them ourselves
 *
 * @param w http.ResponseWriter - used to write error messages to the user if one occurs (if so, will be a 500 error)
 * @param c *Country            - "main struct" passed by reference, as we want to operate on the original
 *                                 rather than a copy of it
 * @param limit int             - how many cities we want to fetch
 *
 * @returns - an empty error if an error occurs, nil otherwise
 */
func FetchCities(w http.ResponseWriter, c *Country, limit int) error {

	// makes a payload, the input in the POST request with the common name we already have
	payload := strings.NewReader("{\"country\": \"" + c.Name.Common + "\"}")

	// makes a POST request with the payload and stores the response + logs potential errors
	resp, errNOW := http.Post(consts.COUNTRIESNOWURL+"countries/cities", "application/json", payload)
	if errNOW != nil {
		log.Println("(FetchCities) Error in POST request: ", errNOW.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return errors.New("")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Error: iso-2 code is not in use. (Error code 3000)", http.StatusNotFound)
	}

	// reads the response body (the json) and stores it as a []byte, so we can unmarshal later + logs potental errors
	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Println("(FetchCities) Error in io.ReadAll: ", errReadAll.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return errors.New("")
	}

	// temporary struct that just contains cities
	var temp JutsCities

	// formats the json we fetched into the temporary struct + logs potential errors
	errJson := json.Unmarshal(body, &temp)
	if errJson != nil {
		log.Println("(FetchCities) There was an error parsing json: ", errJson.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return errors.New("")
	}

	// sorts cities, as they are only partially sorted when we fetch them from CountriesNow
	sort.Strings(c.Cities)

	// appends the first 'limit' elements of the temporary structs slice into the original
	c.Cities = append(c.Cities, temp.Cities[:limit]...)
	return nil
}

/**
 *
 * Prettyprints based on the Country struct
 * In reality we make a similar struct to prettyprint based on how it looks like in the assignment
 * The reason we make a new one and don't use Country, is because it's not correctly formatted
 * This may not be the best workaround, but eh
 *
 * @param w http.ResponseWriter - writes errors to the user if any
 * @param c Country             - the country to be printed
 */
func PrintCountry(w http.ResponseWriter, c Country) {
	var country struct {
		Name       string            `json:"name"`
		Continents []string          `json:"continents"`
		Languages  map[string]string `json:"languages"`
		Population int               `json:"population"`
		Borders    []string          `json:"borders"`
		Flag       string            `json:"flag"`
		Capital    []string          `json:"capital"`
		Cities     []string          `json:"cities"`
	}

	country.Name = c.Name.Common
	country.Continents = c.Continents
	country.Languages = c.Languages
	country.Population = c.Population
	country.Borders = c.Borders
	country.Flag = c.Flag
	country.Capital = c.Capital
	country.Cities = c.Cities

	jsonCOUNTRY, err := json.MarshalIndent(country, "", "    ")
	if err != nil {
		log.Println("(PrintCountry) Could not marshall json: " + err.Error())  // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // :) 500
		return
	}
	fmt.Fprint(w, string(jsonCOUNTRY))
}
