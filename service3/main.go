package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Service 3!")
	})


	log.Printf("Service 3 is running on port 8002\n")
	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		log.Fatal(err)
	}
}
