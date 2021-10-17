package enc

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	c, err := NewConfig()
	if err != nil {
		responseErr(w, 500, "App invalid config")
		return
	}

	err = req.ParseMultipartForm(21 * 1024 * 1024)
	if err != nil {
		responseErr(w, 400, "Parsing multipart form failed")
		return
	}

	form := req.MultipartForm
	mode := ""
	modes, ok := form.Value["mode"]
	if ok {
		mode = modes[0]
	} else {
		mode = "once"
	}

	switch mode {
	case "once":
		handleOnce(w, form, c)
	case "chunked":
		handleChunked(w, form, c)
	default:
		responseErr(w, 400, "Invalid mode param")
	}
}

func handleOnce(w http.ResponseWriter, form *multipart.Form, c Config) {}

func handleChunked(w http.ResponseWriter, form *multipart.Form, c Config) {}

func responseErr(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	_, _ = fmt.Fprintln(w, reason)
	log.Printf("Err: status = %d, reason = %s", status, reason)
}
