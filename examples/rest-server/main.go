package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bygui86/multi-profile"
)

func main() {
	defer profile.CPUProfile(&profile.ProfileConfig{}).Start().Stop()
	defer profile.MemProfile(&profile.ProfileConfig{}).Start().Stop()

	log.Println("Starting handling requests")
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: homePage")
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
