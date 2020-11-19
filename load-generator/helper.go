package load

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func MakeRequest(url string, method string, headers []string, parameters string, timeout int, data string) (int, string) {

	if url != "" {
		url = url + "?" + parameters
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	var req *http.Request
	var err error
	if data != "" {
		jsonData := []byte(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if headers[0] != "" {
		for _, header := range headers {
			splitted := strings.Split(header, ":")
			req.Header.Add(splitted[0], splitted[1])
		}
	}

	resp, err2 := client.Do(req)

	if err2 != nil {
		log.Fatal("Error reading response. ", err)

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	defer resp.Body.Close()
	return resp.StatusCode, string(body)
}
