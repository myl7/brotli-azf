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
	req.Header.Set("content-type", mw.FormDataContentType())
	enc.Handle(w, req)
	res := w.Result()
	_, params, err := mime.ParseMediaType(res.Header.Get("content-type"))
	if err != nil {
		t.Error(err)
		return
	}

	mr := multipart.NewReader(res.Body, params["boundary"])
	foundFile := false
	foundLen := false
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
			return
		}

		key := p.FormName()
		switch key {
		case "file":
			foundFile = true
			output, _ := ioutil.ReadAll(p)

			var b2 bytes.Buffer
			bw := brotli.NewWriter(&b2)
			_, _ = bw.Write(input)
			_ = bw.Close()

			toOutput := b2.Bytes()
			if bytes.Compare(output, toOutput) != 0 {
				t.Errorf("file requires %v but get %v", toOutput, output)
				return
			}
		case "len":
			foundLen = true
			l, _ := ioutil.ReadAll(p)
			if string(l) != "64" {
				t.Errorf("len requires %d but get %s", 64, l)
				return
			}
		}
	}
	if foundFile == false {
		t.Errorf("file not exists")
		return
	}
	if foundLen == false {
		t.Errorf("len not exists")
		return
	}
}
