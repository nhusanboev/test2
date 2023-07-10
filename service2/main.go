package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Service 2!")
	})


	log.Printf("Service 2 is running on port 8001\n")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
