package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(len(r.PathValue("two_letter_country_code")), r.URL.Query().Get("limit"), r.URL.Port())

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

		FetchCountry(w, r, 10) // default limit is 10 (when no limit was explicitly given by the user)

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

		FetchCountry(w, r, limit)

		//
	}

	//

	//

	//

	fmt.Fprint(w, "Success!, nothing to show yet though...")

}

func FetchCountry(w http.ResponseWriter, r *http.Request, limit int) {
	//--
	/*
		resp, err := http.Get(consts.RESTCOUNTRIESURL + "alpha/" + iso)
		if err != nil {
			fmt.Println(err.Error()) // debug
		}
		defer resp.Body.Close()

		body, err1 := io.ReadAll(resp.Body)
		if err1 != nil {
			return
		}

		fmt.Println(string(body))

		jsonStatus, errjson := json.Marshal(string(body))
		if errjson != nil {
			fmt.Printf("Error: %s", errjson.Error())
		}
		fmt.Fprint(w, string(jsonStatus))
	*/
}
