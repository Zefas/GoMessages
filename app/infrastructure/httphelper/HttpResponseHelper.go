package httphelper
import (
	"net/http"
	"io/ioutil"
)


func ReadBodyAsString(req *http.Request) (string, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}