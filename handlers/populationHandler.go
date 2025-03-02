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

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	iso := r.PathValue("two_letter_country_code") // get iso code
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 must be a 2 letter code. (Error code 200)", http.StatusBadRequest) // 400
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
			fmt.Fprintln(w, "Error when fetching population: "+errPopNoArgs.Error()) // TODO
			return
		}

	} else {
		timeframe := strings.Split(limit, "-")
		fmt.Fprintln(w, "Limit with args: \""+limit+"\"       timeframe:", timeframe)
		if len(timeframe) != 2 {
			http.Error(w, "Expected 2 arguments, got "+fmt.Sprint(len(timeframe)), http.StatusBadRequest) // TODO
			return                                                                                        // need 2 args
		}
		if timeframe[0] == "" || timeframe[1] == "" {
			http.Error(w, "One or more arguments are empty", http.StatusBadRequest) // TODO
			return                                                                  // one or more empty args
		}

		//

		//

		//

		errPopWithArgs := FetchPopulation(w, iso3, timeframe[0], timeframe[1])
		if errPopWithArgs != nil {
			fmt.Fprintln(w, "Error when fetching population: "+errPopWithArgs.Error()) // TODO
		}
	}

}

func GetCountry(w http.ResponseWriter, iso string) (string, error) {
	var country struct {
		Iso3 string `json:"cca3"`
	}

	resp, errGet := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + "?fields=cca3")
	if errGet != nil {
		log.Println("(FetchCountry) Error in http.Get: ", errGet.Error())      // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // :) 500
		return "", errors.New("")                                              // :) ???
	}
	defer resp.Body.Close()

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		log.Println("(FetchCountry) Error in io.ReadAll: ", errReadAll.Error()) // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError)  // :) 500
		return "", errors.New("")                                               // :) ???
	}

	errJson := json.Unmarshal(body, &country)
	if errJson != nil {
		log.Println("(FetchCountry) Error parsing json with json.Unmarshal: ", errJson.Error()) // :)
		http.Error(w, "Internal server error", http.StatusInternalServerError)                  // :) 500
		return "", errors.New("")                                                               // :) ???
	}

	if country.Iso3 == "" {
		http.Error(w, "Could not retrieve an iso3 code from iso2 code \""+iso+"\". (Error code 201)", http.StatusNotFound) // :) 404
		return "", errors.New("")                                                                                          // :) ???
	}

	return country.Iso3, nil
}

//

//

//

//

func FetchPopulation(w http.ResponseWriter, iso3, min, max string) error {
	var start, end int

	if min != "" {
		s, errConvStart := strconv.Atoi(min)
		if errConvStart != nil {
			return errors.New("start year must be a number") // TODO
		}
		start = s
	} else {
		start = 0
	}
	if max != "" {
		e, errConvEnd := strconv.Atoi(max)
		if errConvEnd != nil {
			return errors.New("end year must be a number") // TODO
		}
		end = e
	} else {
		end = time.Now().Year()
	}
	if start > end {
		return errors.New("start year is greater than end year") // TODO
	}

	var wrapper struct {
		Mean int `json:"mean"`
		Data struct {
			PopulationCounts []struct {
				Year  int `json:"year"`
				Value int `json:"value"`
			} `json:"populationCounts"`
		} `json:"data"`
	}

	payload := strings.NewReader(`{"iso3": "` + iso3 + `"}`)

	// Makes a post request and handles errors. Defers closing the body response.
	resp, errPost := http.Post(consts.COUNTRIESNOWURL+"countries/population", "application/json", payload)

	if errPost != nil {
		log.Print("(FetchPopulation) Error in post request: ", errPost.Error())                  // TODO
		return errors.New(fmt.Sprint(http.StatusInternalServerError) + " internal server error") // TODO
	}
	defer resp.Body.Close()

	fmt.Println("Response: ", resp) // DELETEME

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		return errors.New("error in io.ReadAll: " + errReadAll.Error()) // TODO
	}

	errJson := json.Unmarshal(body, &wrapper)
	if errJson != nil {
		fmt.Println("(FetchCities) There was an error parsing json: ", errJson.Error()) // TODO
	}

	var i, j = 0, 0

	// finds first instance that matches
	for ; i < len(wrapper.Data.PopulationCounts); i++ {
		if start <= wrapper.Data.PopulationCounts[i].Year {
			break
		}
	}
	for j = len(wrapper.Data.PopulationCounts) - 1; 0 < j; j-- {
		if end+1 >= wrapper.Data.PopulationCounts[j].Year { // +1 to include the end year
			break
		}
	}

	wrapper.Data.PopulationCounts = wrapper.Data.PopulationCounts[i:j]

	// calculates sum of all years. 'val' is a struct, which is why we do val.Value
	var sum = 0
	for _, val := range wrapper.Data.PopulationCounts {
		sum += val.Value
	}
	wrapper.Mean = sum / len(wrapper.Data.PopulationCounts)

	fmt.Print("\n\n\n\n\n\n", wrapper.Mean, "\n\n\n\n\n\n") // DELETEME

	//

	//
	// vvv DEBUG vvv
	jsonStatus, errjson := json.MarshalIndent(wrapper, "", "    ")
	if errjson != nil {
		fmt.Println("Error: ", errjson.Error()) // TODO
	}
	fmt.Fprintln(w, string(jsonStatus))

	return nil
}
