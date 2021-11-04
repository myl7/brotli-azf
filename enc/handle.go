package enc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"mime/multipart"
	"net/http"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	c, err := NewConfig()
	if err != nil {
		reportErr(w, 500, "app invalid config", err, nil)
		return
	}

	mr, err := req.MultipartReader()
	if err != nil {
		reportErr(w, 400, "getting multipart form failed", err, nil)
		return
	}

	params, ok := ParseParams(w, mr, c)
	if !ok {
		return
	}

	switch params["mode"].(string) {
	case "once":
		handleOnce(w, params, c)
	// case "chunked":
	// 	handleChunked(w, params, c)
	default:
		reportErr(w, 400, "invalid mode param", nil, nil)
	}
}

func handleOnce(w http.ResponseWriter, params map[string]interface{}, _ Config) {
	res := multipart.NewWriter(w)
	defer func(res *multipart.Writer) {
		err := res.Close()
		if err != nil {
			reportErr(w, 500, "writing failed", err, nil)
		}
	}(res)
	w.Header().Set("content-type", res.FormDataContentType())

	encw, err := res.CreateFormFile("file", "bin")
	if err != nil {
		reportErr(w, 500, "creating form failed", err, nil)
		return
	}

	enc := brotli.NewWriter(encw)
	_, err = bytes.NewReader(params["file"].([]byte)).WriteTo(enc)
	if err != nil {
		reportErr(w, 400, "brotli decoding failed", err, nil)
		return
	}
	err = enc.Close()
	if err != nil {
		reportErr(w, 400, "brotli decoding failed", err, nil)
		return
	}

	lenw, err := res.CreateFormField("len")
	if err != nil {
		reportErr(w, 500, "creating form failed", err, nil)
		return
	}

	_, err = lenw.Write([]byte(fmt.Sprintf("%d", params["len"])))
	if err != nil {
		reportErr(w, 500, "writing failed", err, nil)
		return
	}

	reportOnceOk(params["len"].(int))
}

func handleChunked(w http.ResponseWriter, params map[string]interface{}, c Config) {}

var errFailed = errors.New("failed")
