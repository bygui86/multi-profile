package main

import (
	"fmt"
	"net/http"

	"github.com/bygui86/multi-profile"

	"github.com/bygui86/multi-profile/examples/custom-logger/logging"
)

func main() {
	logging.InitGlobalLogger()

	defer profile.CPUProfile(&profile.ProfileConfig{Path: "/multi-profile", Logger: logging.SugaredLog}).Start().Stop()
	defer profile.MemProfile(&profile.ProfileConfig{Path: "/multi-profile", Logger: logging.SugaredLog}).Start().Stop()

	logging.Log.Info("Starting handling requests")
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("handle homePage endpoint")

	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	logging.SugaredLog.Fatal(http.ListenAndServe(":8080", nil))
}
