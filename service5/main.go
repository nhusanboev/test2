package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Service 5!")
	})


	log.Printf("Service 5 is running on port 8004\n")
	err := http.ListenAndServe(":8004", nil)
	if err != nil {
		log.Fatal(err)
	}
}
