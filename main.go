package main

import (
	"brotli-azf/enc"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if ok {
		port = val
	}

	http.HandleFunc("/api/enc", enc.Handle)

	listen := fmt.Sprintf("Listening on :%s", port)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
