package api

import (
	"fmt"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/", homePageHandler)

	fmt.Println("API listening on port 3001")
	log.Panic(http.ListenAndServe(":3001", nil))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello world")
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
