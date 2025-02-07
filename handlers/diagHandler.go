package handlers

import (
	"fmt"
	"log"
	"net/http"
)

const LINEBREAK = "\n"

func DiagHandler(w http.ResponseWriter, r *http.Request) {

	output := "Request:" + LINEBREAK
	output += "URL Path" + r.URL.Path + LINEBREAK
	output += "Method" + r.Method + LINEBREAK

	output += LINEBREAK + "Headers: " + LINEBREAK

	for k, v := range r.Header {
		for _, vv := range v {
			output += k + ": " + vv + LINEBREAK
		}
	}

	_, err := fmt.Fprint(w, "%V", output)
	if err != nil {
		log.Print("An error occured: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
