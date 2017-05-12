package nightscout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	siteEnvVar      = "NIGHTSCOUT_SITE"
	apiSecretEnvVar = "NIGHTSCOUT_API_SECRET"
)

var (
	verbose  = false
	noUpload = false
)

func Verbose() bool {
	return verbose
}

func SetVerbose(flag bool) {
	verbose = flag
}

func NoUpload() bool {
	return noUpload
}

func SetNoUpload(flag bool) {
	noUpload = flag
}

func Site() (string, error) {
	site := os.Getenv(siteEnvVar)
	if len(site) == 0 {
		return "", fmt.Errorf("%s is not set", siteEnvVar)
	}
	return site, nil
}

func ApiSecret() (string, error) {
	secret := os.Getenv(apiSecretEnvVar)
	if len(secret) == 0 {
		return "", fmt.Errorf("%s is not set", apiSecretEnvVar)
	}
	return secret, nil
}

func RestOperation(op string, api string, v interface{}) ([]byte, error) {
	req, err := makeRequest(op, api, v)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return []byte("[]"), nil
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if verbose && err == nil {
		log.Print(indentJson(result))
	}
	return result, err
}

func makeRequest(op string, api string, v interface{}) (*http.Request, error) {
	u, err := makeURL(op, api)
	if err != nil {
		return nil, err
	}
	if verbose || noUpload {
		log.Printf("%s %v", op, u)
		if v != nil {
			log.Print(Json(v))
		}
	}
	if noUpload && op != "GET" {
		return nil, nil
	}
	r, err := makeReader(v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(op, u, r)
	if err != nil {
		return nil, err
	}
	err = addHeaders(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func makeURL(op string, api string) (string, error) {
	site, err := Site()
	if err != nil {
		return "", err
	}
	u, err := url.Parse(site + "/api/v1/" + api)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func makeReader(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func addHeaders(req *http.Request) error {
	secret, err := ApiSecret()
	if err != nil {
		return err
	}
	req.Header.Add("api-secret", secret)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	return nil
}

func Get(api string, result interface{}) error {
	data, err := RestOperation("GET", api, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func Upload(op string, api string, param interface{}) error {
	if op == "GET" {
		log.Panicf("GET used for upload to %s", api)
	}
	_, err := RestOperation(op, api, param)
	return err
}

func indentJson(data []byte) string {
	buf := bytes.Buffer{}
	err := json.Indent(&buf, data, "", "  ")
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func Hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

func Username() string {
	u := os.Getenv("USER")
	if len(u) == 0 {
		return "unknown"
	}
	return u
}

func Json(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return indentJson(data)
}
