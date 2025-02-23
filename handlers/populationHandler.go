package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	iso := r.PathValue("two_letter_country_code")
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 can only be a 2 letter code. Error code "+fmt.Sprint(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	country, err := GetCountry(w, iso)
	if err != nil {
		fmt.Fprint(w, err.Error()) // might change ???
	}

	//

	fmt.Println("Country retrived: ", country)
	fmt.Fprint(w, country)
}

func GetCountry(w http.ResponseWriter, iso string) (string, error) {
	var country struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
	}

	resp, errGet := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + "?fields=name")
	if errGet != nil {
		fmt.Println("(FetchCountry) Error in http.Get: ", errGet.Error()) // debug
		return "", errors.New("Error in http.Get: " + errGet.Error())     // change me prolly
	}
	defer resp.Body.Close()

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		fmt.Println("(FetchCountry) Error in io.ReadAll: ", errReadAll.Error()) // debug
		return "", errors.New("Error in io.ReadAll: " + errReadAll.Error())     // change me prolly
	}

	errJson := json.Unmarshal(body, &country)
	if errJson != nil {
		fmt.Println("(FetchCountry) There was an error parsing json: ", errJson.Error())
	}

	if country.Name.Common == "" {
		return "", errors.New("Could not retrieve a country from iso code " + iso)
	}

	return country.Name.Common, nil
}
