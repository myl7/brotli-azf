package test

import (
	"brotli-azf/enc"
	"bytes"
	"github.com/andybalholm/brotli"
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func TestEncNoParamWork(t *testing.T) {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	f, _ := mw.CreateFormFile("file", "bin")

	input := make([]byte, 64)
	rand.Read(input)
	_, _ = bytes.NewReader(input).WriteTo(f)

	_ = mw.Close()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://example.com/enc", &body)
	enc.Handle(w, req)
	res := w.Result()
	_, params, err := mime.ParseMediaType(res.Header.Get("content-type"))
	if err != nil {
		t.Error(err)
		return
	}

	l, ok := params["len"]
	if !ok {
		t.Errorf("len not exists")
		return
	} else if l != "64" {
		t.Errorf("len requres %d but gets %s", 64, l)
		return
	}

	mr := multipart.NewReader(res.Body, params["boundary"])
	foundFile := false
	for p, err := mr.NextPart(); err != io.EOF; {
		if err != nil {
			t.Error(err)
			return
		}

		if p.FormName() == "file" {
			foundFile = true
			output, _ := ioutil.ReadAll(p)

			var b2 bytes.Buffer
			bw := brotli.NewWriter(&b2)
			_, _ = bw.Write(input)

			toOutput := b2.Bytes()
			if bytes.Compare(output, toOutput) != 0 {
				t.Errorf("file requires %v but get %v", toOutput, output)
				return
			}
		}
	}
	if foundFile == false {
		t.Errorf("file not exists")
		return
	}
}
