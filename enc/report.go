package enc

import (
	"log"
	"net/http"
)

func reportErr(w http.ResponseWriter, status int, reason string, err error) {
	http.Error(w, reason, status)
	if err != nil {
		log.Printf("err: status = %d; reason = %s; original err = %s", status, reason, err.Error())
	} else {
		log.Printf("err: status = %d; reason = %s", status, reason)
	}
}

func reportOnceOk(inputLen int) {
	log.Printf("ok: input len = %d", inputLen)
}
