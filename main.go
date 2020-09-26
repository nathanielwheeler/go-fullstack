package main

import (
	"fmt"
	"net/http"

  "github.com/gorilla/mux"
  
  "github.com/nathanielwheeler/go-fullstack/api/controllers"
  "github.com/nathanielwheeler/go-fullstack/api/models"
)

func main() {
  // Initialize services
  services, err := services.NewValuesService

	// Router Initialization
  r := mux.NewRouter()
  
  // Initialize controllers
  valuesC := controllers.NewValues(services.Value)

	// Asset Routes
	assetHandler := http.FileServer(http.Dir("./client/assets/"))
	r.PathPrefix("/css/").Handler(assetHandler)
	apphandler := http.FileServer(http.Dir("./client/app/js"))
  r.PathPrefix("/js/").Handler(apphandler)

	// Value Routes
	r.HandleFunc("/", valuesC.Index).Methods("GET")

	// Start server
	port := fmt.Sprintf(":%d", 6789)
	fmt.Printf("Now listening on %s...\n", port)
	http.ListenAndServe(port, r)
}
