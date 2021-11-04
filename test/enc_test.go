package test

import (
	"brotli-azf/enc"
	"bytes"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"testing"
)

func TestEncNoParamWork(t *testing.T) {
	input := make([]byte, 64)
	rand.Read(input)
	body, contentType := CreateForm(bytes.NewReader(input), nil, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://example.com/enc", body)
	req.Header.Set("content-type", contentType)
	enc.Handle(w, req)
	res := w.Result()
	output, _ := ioutil.ReadAll(res.Body)

	var b2 bytes.Buffer
	bw := brotli.NewWriter(&b2)
	_, _ = bw.Write(input)
	_ = bw.Close()
	expected := b2.Bytes()
	if bytes.Compare(output, expected) != 0 {
		t.Errorf("file requires %v but get %v", expected, output)
		return
	}
}
