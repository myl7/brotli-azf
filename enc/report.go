package enc

import (
	"log"
	"net/http"
)

func reportErr(w http.ResponseWriter, status int, reason string, err error, extra interface{}) {
	http.Error(w, reason, status)
	if err != nil {
		log.Printf("err: status = %d; reason = %s; original err = %s; extra = %v", status, reason, err.Error(), extra)
	} else {
		log.Printf("err: status = %d; reason = %s; extra = %v", status, reason, extra)
	}
}

func reportOk(inputLen int64) {
	log.Printf("ok: input len = %d", inputLen)
}
