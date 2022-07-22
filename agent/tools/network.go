package tools

import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err = ioutil.ReadAll(resp.Body)
	return
}
