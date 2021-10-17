package enc

import (
	"fmt"
	"mime/multipart"
	"strconv"
)

func getRangedIntParam(form *multipart.Form, name string, defaultVal int, min int, max int) (int, error) {
	ss, ok := form.Value[name]
	if ok {
		s := ss[0]
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: not number", name)
		}

		if !(min <= d && d <= max) {
			return 0, fmt.Errorf("invalid %s: out of range", name)
		}

		return d, nil
	} else {
		return defaultVal, nil
	}
}

func getQuality(form *multipart.Form) (int, error) {
	return getRangedIntParam(form, "quality", 11, 0, 11)
}

func getLGWin(form *multipart.Form) (int, error) {
	return getRangedIntParam(form, "lgwin", 0, 10, 24)
}
