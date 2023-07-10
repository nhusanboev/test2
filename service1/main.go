package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Service 1!")
	})


	log.Printf("Service 1 is running on port 8000\n")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
