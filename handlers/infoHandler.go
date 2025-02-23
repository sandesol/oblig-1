package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	Cities     []string          `json:"cities"`
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(len(r.PathValue("two_letter_country_code")), r.URL.Query().Get("limit"), r.URL.Port())

	country := Country{}
	iso := r.PathValue("two_letter_country_code")
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 can only be a 2 letter code. Error code "+fmt.Sprint(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("limit")
	if query == "" {
		fmt.Println("No limit :)") // DELETEME

		//

		//

		//

		//

		FetchCountry(w, r, 10, &country, iso) // default limit is 10 (when no limit was explicitly given by the user)

		//
	} else {
		fmt.Println("A limit have been set :))))") // DELETEME

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Error: limit must be an integer. Error code "+fmt.Sprint(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if limit < 0 {
			http.Error(w, "Error: limit must be a positive integer. Error code "+fmt.Sprint(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//

		//

		//

		//

		FetchCountry(w, r, limit, &country, iso)

		//
	}

	errprint := PrintCountry(w, country)
	if errprint != nil {
		fmt.Println(errprint.Error())
	}

	//

	//

	//

	fmt.Fprint(w, "Success!, nothing to show yet though...") // DELETEME

}

/**
 *	For some odd reason, r.PathValue can only be read once, which means we have to pass iso as a parameter
 *
 *
 */
func FetchCountry(w http.ResponseWriter, r *http.Request, limit int, country *Country, iso string) {

	resp, err := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + consts.RESTCOUNTRIESFILTER)
	if err != nil {
		fmt.Println(err.Error()) // debug
	}
	defer resp.Body.Close()

	body, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		return
	}

	err2 := json.Unmarshal(body, country)
	if err2 != nil {
		fmt.Println("There was an error parsing json: ", err2.Error())
	}
}

func PrintCountry(w http.ResponseWriter, c Country) error {
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
		return errors.New("Could not marshall json: " + err.Error())
	}
	fmt.Fprint(w, string(jsonCOUNTRY))

	return nil
}
