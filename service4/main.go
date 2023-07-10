package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Service 4!")
	})


	log.Printf("Service 4 is running on port 8003\n")
	err := http.ListenAndServe(":8003", nil)
	if err != nil {
		log.Fatal(err)
	}
}
