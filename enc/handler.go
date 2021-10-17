package enc

import (
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	c, err := NewConfig()
	if err != nil {
		responseErr(w, 500, "app invalid config")
		return
	}

	err = req.ParseMultipartForm(1*1024*1024 + int64(config.InputMaxsize))
	if err != nil {
		responseErr(w, 400, "parsing multipart form failed")
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
		responseErr(w, 400, "invalid mode param")
	}
}

func handleOnce(w http.ResponseWriter, form *multipart.Form, c Config) {
	f, err := getFile(form)
	if err != nil {
		responseErr(w, 400, err.Error())
		return
	}

	defer func(f multipart.File) {
		_ = f.Close()
	}(f)

	res := multipart.NewWriter(w)
	defer func(res *multipart.Writer) {
		_ = res.Close()
	}(res)

	encw, err := res.CreateFormFile("file", "bin")
	if err != nil {
		responseErr(w, 500, "creating form failed")
		return
	}

	enc := brotli.NewWriter(encw)
	defer func(enc *brotli.Writer) {
		err := enc.Close()
		if err != nil {
			responseErr(w, 400, "brotli decoding failed")
		}
	}(enc)

	n, err := io.CopyN(enc, f, int64(c.InputMaxsize))
	if err != nil {
		responseErr(w, 400, "brotli decoding failed")
		return
	}

	lenw, err := res.CreateFormField("len")
	if err != nil {
		responseErr(w, 500, "creating form failed")
		return
	}

	_, err = lenw.Write([]byte(fmt.Sprintf("%d", n)))
	if err != nil {
		responseErr(w, 500, "creating form failed")
		return
	}

	w.Header().Set("content-type", res.FormDataContentType())
}

func handleChunked(w http.ResponseWriter, form *multipart.Form, c Config) {
	f, err := getFile(form)
	if err != nil {
		responseErr(w, 400, err.Error())
		return
	}

	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
}

func responseErr(w http.ResponseWriter, status int, reason string) {
	http.Error(w, reason, status)
	log.Printf("error status %d reason %s", status, reason)
}

func getFile(form *multipart.Form) (multipart.File, error) {
	files, ok := form.File["file"]
	if !ok {
		return nil, errors.New("no file param")
	}

	file := files[0]
	f, err := file.Open()
	if err != nil {
		return nil, errors.New("opening file failed")
	}

	return f, nil
}
