package enc

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

func ParseParams(w http.ResponseWriter, mr *multipart.Reader, c Config) (map[string]interface{}, bool) {
	res := make(map[string]interface{})
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			reportErr(w, 400, "invalid multipart form", err, nil)
			return nil, false
		}

		key := p.FormName()
		if key == "file" {
			res["file"] = p
			break
		}

		b, err := ioutil.ReadAll(p)
		if err != nil {
			reportErr(w, 500, "reading failed", err, nil)
			return nil, false
		}
		val := string(b)

		switch key {
		case "quality":
			d, err := checkRangedIntParam(key, val, 0, 11)
			if err != nil {
				return nil, false
			}

			res[key] = d
		case "lgwin":
			d, err := checkRangedIntParam(key, val, 10, 24)
			if err != nil {
				return nil, false
			}

			res[key] = d
		}
	}

	_, ok := res["file"]
	if !ok {
		reportErr(w, 400, "no file to be processed", nil, nil)
		return nil, false
	}

	_, ok = res["mode"]
	if !ok {
		res["mode"] = "once"
	}

	_, ok = res["quality"]
	if !ok {
		res["quality"] = 11
	}

	_, ok = res["lgwin"]
	if !ok {
		res["lgwin"] = 0
	}

	return res, true
}

func checkRangedIntParam(name string, val string, min int, max int) (int, error) {
	d, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: not number", name)
	}

	if !(min <= d && d <= max) {
		return 0, fmt.Errorf("invalid %s: out of range", name)
	}

	return d, nil
}
