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

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	country := Country{}
	iso := r.PathValue("two_letter_country_code")
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 must be a 2 letter code. (Error code 100)", http.StatusBadRequest)
		return
	}

	FetchCountry(w, &country, iso) // fetches countries first, so we can use country name as a parameter in the post request in FetchCities()

	query := r.URL.Query().Get("limit")
	if query == "" {

		FetchCities(w, &country, 10) // No limit set, defaults to 10

	} else {

		limit, err := strconv.Atoi(r.URL.Query().Get("limit")) // gets limit, converts it to int

		if err != nil {
			http.Error(w, "Error: limit must be an integer. (Error code 101)", http.StatusBadRequest)
			return
		}
		if limit < 0 {
			http.Error(w, "Error: limit must be a positive integer. (Error code 102)", http.StatusBadRequest)
			return
		}

		FetchCities(w, &country, limit) // passed both tests, limit is valid and we can fetch cities

	}

	PrintCountry(w, country) // prettyprints the complete country
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
 * Fetches all cities via countriesNow
 *
 * Country is passed by reference
 */
func FetchCities(w http.ResponseWriter, c *Country, limit int) {

	payload := strings.NewReader("{\"country\": \"" + c.Name.Common + "\"}")
	fmt.Println(c.Name.Common)

	resp, errNOW := http.Post(consts.COUNTRIESNOWURL+"countries/cities", "application/json", payload)
	if errNOW != nil {
		log.Println("(FetchCities) Error in POST request: ", errNOW.Error())   // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // :) 500
		return
	}
	defer resp.Body.Close()

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Println("(FetchCities) Error in io.ReadAll: ", errReadAll.Error()) // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // :) 500
		return
	}

	// we use a temporary struct to wrap cities in, because it's less annoying to deal with than working with the original
	var temp struct {
		Cities []string `json:"data"`
	}

	errJson := json.Unmarshal(body, &temp)
	if errJson != nil {
		log.Println("(FetchCities) There was an error parsing json: ", errJson.Error()) // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError)          // :) 500
	}

	// appends the first 'limit' elements of the temporary structs slice into the original
	c.Cities = append(c.Cities, temp.Cities[:limit]...)
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
