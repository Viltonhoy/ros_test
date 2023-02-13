package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", helloDocker)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func helloDocker(w http.ResponseWriter, r *http.Request) {
	var s = "Hello Docker!"

	if ev := os.Getenv("VALUE"); ev != "" {
		s = ev
	}
	fmt.Fprint(w, s)
}
