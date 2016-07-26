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

func Verbose() bool        { return verbose }
func SetVerbose(flag bool) { verbose = flag }

func NoUpload() bool        { return noUpload }
func SetNoUpload(flag bool) { noUpload = flag }

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
	site, err := Site()
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(site + "/api/v1/" + api)
	if err != nil {
		return nil, err
	}
	data := []byte{}
	if v != nil {
		data, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	if verbose || noUpload {
		log.Printf("%s %v", op, u)
		if v != nil {
			log.Printf("%s", string(data))
		}
	}
	if noUpload && op != "GET" {
		return []byte{'[', ']'}, nil
	}
	reader := io.Reader(nil)
	if len(data) != 0 {
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(op, u.String(), reader)
	if err != nil {
		return nil, err
	}
	secret, err := ApiSecret()
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-secret", secret)
	req.Header.Add("accept", "application/json")
	if len(data) != 0 {
		req.Header.Add("content-type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if verbose && err == nil {
		log.Println(indentJson(result))
	}
	return result, err
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
