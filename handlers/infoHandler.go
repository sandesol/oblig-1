package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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
 *
 * A struct that just contains cities.
 * This struct exists solely to make filtering operations more convenient.
 *
 */
type JutsCities struct {
	Cities []string `json:"data"`
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	country := Country{}
	iso := r.PathValue("two_letter_country_code")
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 must be a 2 letter code. (Error code 100)", http.StatusBadRequest)
		return
	}

	// fetches countries first, so we can use country name as a parameter in the post request in FetchCities()
	FetchCountry(w, &country, iso)

	query := r.URL.Query().Get("limit")
	if query == "" {

		// No limit set, defaults to 10
		FetchCities(w, &country, 10)

	} else {

		// gets limit, converts it to int
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil {
			http.Error(w, "Error: limit must be an integer. (Error code 101)", http.StatusBadRequest)
			return
		}
		if limit < 0 {
			http.Error(w, "Error: limit must be a positive integer. (Error code 102)", http.StatusBadRequest)
			return
		}

		// passed both tests, limit is valid and we can fetch cities
		FetchCities(w, &country, limit)

	}

	// prettyprints the complete country
	PrintCountry(w, country)
}

/**
 *	For some odd reason, r.PathValue can only be read once, which means we have to pass iso as a parameter
 *
 *	Country is passed by reference
 */
func FetchCountry(w http.ResponseWriter, c *Country, iso string) {

	resp, errGet := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + consts.RESTCOUNTRIESFILTER)
	if errGet != nil {
		fmt.Println("(FetchCountry) Error in http.Get: ", errGet.Error()) // debug
		return
	}
	defer resp.Body.Close()

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		fmt.Println("(FetchCountry) Error in io.ReadAll: ", errReadAll.Error()) // debug
		return
	}

	errJson := json.Unmarshal(body, c)
	if errJson != nil {
		fmt.Println("(FetchCountry) There was an error parsing json: ", errJson.Error())
	}
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
 */
func FetchCities(w http.ResponseWriter, c *Country, limit int) {

	// makes a payload, the input in the POST request with the common name we already have
	payload := strings.NewReader("{\"country\": \"" + c.Name.Common + "\"}")

	// makes a POST request with the payload and stores the response + logs potential errors
	resp, errNOW := http.Post(consts.COUNTRIESNOWURL+"countries/cities", "application/json", payload)
	if errNOW != nil {
		log.Println("(FetchCities) Error in POST request: ", errNOW.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return
	}
	defer resp.Body.Close()

	// reads the response body (the json) and stores it as a []byte, so we can unmarshal later + logs potental errors
	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Println("(FetchCities) Error in io.ReadAll: ", errReadAll.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
		return
	}

	// temporary struct that just contains cities
	var temp JutsCities

	// formats the json we fetched into the temporary struct + logs potential errors
	errJson := json.Unmarshal(body, &temp)
	if errJson != nil {
		log.Println("(FetchCities) There was an error parsing json: ", errJson.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
	}

	// appends the first 'limit' elements of the temporary structs slice into the original
	c.Cities = append(c.Cities, temp.Cities[:limit]...)

	// sorts cities, as they are only partially sorted when we fetch them from CountriesNow
	sort.Strings(c.Cities)
}

/**
 *
 * Prettyprints based on the Country struct
 * In reality we make a similar struct to prettyprint based on how it looks like in the assignment
 * The reason we make a new one and don't use Country, is because it's not correctly formatted
 * This may not be the best workaround, but eh
 *
 * Returns an error if it json.MarshallIndent failed, and nil otherwise
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
