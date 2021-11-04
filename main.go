package main

import (
	"fmt"
	"github.com/myl7/brotli-azf/enc"
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

	listen := fmt.Sprintf(":%s", port)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
