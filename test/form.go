package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

func CreateForm(file io.Reader, quality *int, lgwin *int) (io.Reader, string) {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	f, _ := mw.CreateFormFile("file", "bin")
	_, _ = bufio.NewReader(file).WriteTo(f)

	if quality != nil {
		f, _ := mw.CreateFormField("quality")
		_, _ = f.Write([]byte(fmt.Sprintf("%d", *quality)))
	}

	if lgwin != nil {
		f, _ := mw.CreateFormField("lgwin")
		_, _ = f.Write([]byte(fmt.Sprintf("%d", *lgwin)))
	}

	_ = mw.Close()
	return &body, mw.FormDataContentType()
}
