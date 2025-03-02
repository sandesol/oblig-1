package handlers

import (
	"assignment1/consts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	iso := r.PathValue("two_letter_country_code") // get iso code
	if len(iso) != 2 {
		http.Error(w, "Error: iso-2 can only be a 2 letter code. Error code "+fmt.Sprint(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	iso3, err := GetCountry(w, iso) // get country name
	if err != nil {
		fmt.Fprintln(w, err.Error()) // might change ???
		return
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" { // no args, use whole range

		fmt.Fprintln(w, "Limit with no args: \""+limit+"\"")

		errPopNoArgs := FetchPopulation(w, iso3, "", "")
		if errPopNoArgs != nil {
			fmt.Fprintln(w, "Error when fetching population: "+errPopNoArgs.Error()) // --
			return
		}

	} else {
		timeframe := strings.Split(limit, "-")
		fmt.Fprintln(w, "Limit with args: \""+limit+"\"       timeframe:", timeframe)
		if len(timeframe) != 2 {
			http.Error(w, "Expected 2 arguments, got "+fmt.Sprint(len(timeframe)), http.StatusBadRequest) // --
			return                                                                                        // need 2 args
		}
		if timeframe[0] == "" || timeframe[1] == "" {
			http.Error(w, "One or more arguments are empty", http.StatusBadRequest) // --
			return                                                                  // one or more empty args
		}

		//

		//

		//

		errPopWithArgs := FetchPopulation(w, iso3, timeframe[0], timeframe[1])
		if errPopWithArgs != nil {
			fmt.Fprintln(w, "Error when fetching population: "+errPopWithArgs.Error()) // --
		}
	}

}

func GetCountry(w http.ResponseWriter, iso string) (string, error) {
	var country struct {
		Iso3 string `json:"cca3"`
	}

	resp, errGet := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso + "?fields=cca3")
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

	if country.Iso3 == "" {
		return "", errors.New("Could not retrieve an iso3 code from iso2 code \"" + iso + "\"")
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
			return errors.New("start year must be a number")
		}
		start = s
	} else {
		start = 0
	}
	if max != "" {
		e, errConvEnd := strconv.Atoi(max)
		if errConvEnd != nil {
			return errors.New("end year must be a number")
		}
		end = e
	} else {
		end = time.Now().Year()
	}
	if start > end {
		return errors.New("start year is greater than end year")
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
		return errors.New("Error in post request: " + errPost.Error()) // --
	}
	defer resp.Body.Close()

	fmt.Println("Payload: ", payload, "     ", "URL: ", consts.RESTCOUNTRIESURL+"countries/population") /////////////////////////
	fmt.Println("Response: ", resp)

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		return errors.New("error in io.ReadAll: " + errReadAll.Error())
	}

	errJson := json.Unmarshal(body, &wrapper)
	if errJson != nil {
		fmt.Println("(FetchCities) There was an error parsing json: ", errJson.Error())
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

	fmt.Print("\n\n\n\n\n\n", wrapper.Mean, "\n\n\n\n\n\n")

	//

	//
	// vvv DEBUG vvv
	jsonStatus, errjson := json.MarshalIndent(wrapper, "", "    ")
	if errjson != nil {
		fmt.Println("Error: ", errjson.Error())
	}
	fmt.Fprintln(w, string(jsonStatus))

	return nil
}
