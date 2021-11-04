package enc

import (
	"bufio"
	"github.com/andybalholm/brotli"
	"io"
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

	w.Header().Set("content-type", "application/octet-stream")

	switch params["mode"].(string) {
	case "stream":
		processStream(w, params, c)
	case "chunked":
		processChunked(w, params, c)
	default:
		reportErr(w, 400, "unknown mode", nil, nil)
	}
}

func processStream(w http.ResponseWriter, params map[string]interface{}, _ Config) {
	enc := brotli.NewWriterOptions(w, brotli.WriterOptions{
		Quality: params["quality"].(int),
		LGWin:   params["lgwin"].(int),
	})
	file := bufio.NewReader(params["file"].(io.Reader))
	l, err := file.WriteTo(enc)
	if err != nil {
		reportErr(w, 400, "brotli decoding failed", err, nil)
		return
	}

	err = enc.Close()
	if err != nil {
		reportErr(w, 400, "brotli decoding failed", err, nil)
		return
	}

	reportOk(l)
}

func processChunked(w http.ResponseWriter, params map[string]interface{}, _ Config) {}
