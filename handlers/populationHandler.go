package handlers

import (
	"fmt"
	"net/http"
)

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "-Population data goes here-")
	fmt.Fprintln(w, r)
}
