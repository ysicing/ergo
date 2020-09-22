package githubapi

import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
