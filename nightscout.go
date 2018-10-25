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
	"path"
	"regexp"
	"strings"
)

const (
	siteEnvVar      = "NIGHTSCOUT_SITE"
	apiSecretEnvVar = "NIGHTSCOUT_API_SECRET"
	deviceEnvVar    = "NIGHTSCOUT_DEVICE"
)

var (
	verbose  = false
	noUpload = false
)

// Verbose returns the value of the verbose flag.
func Verbose() bool {
	return verbose
}

// SetVerbose sets the value of the verbose flag.
func SetVerbose(flag bool) {
	verbose = flag
}

// NoUpload returns the value of the noUpload flag.
func NoUpload() bool {
	return noUpload
}

// SetNoUpload sets the value of the noUpload flag.
func SetNoUpload(flag bool) {
	noUpload = flag
}

func sitename() (string, error) {
	site := os.Getenv(siteEnvVar)
	if len(site) == 0 {
		return "", fmt.Errorf("%s is not set", siteEnvVar)
	}
	return site, nil
}

func apiSecret() (string, error) {
	secret := os.Getenv(apiSecretEnvVar)
	if len(secret) == 0 {
		return "", fmt.Errorf("%s is not set", apiSecretEnvVar)
	}
	return secret, nil
}

// RestOperation performs an operation on a Nightscout API,
// with optional JSON data, and returns the result.
func RestOperation(op string, api string, v interface{}) ([]byte, error) {
	if op == "GET" && v != nil {
		log.Panicf("GET %s operation with data", api)
	}
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
	result, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if verbose && err == nil {
		log.Print(indentJSON(result))
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
			log.Print(JSON(v))
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
	site, err := sitename()
	if err != nil {
		return "", err
	}
	siteURL, err := url.Parse(site)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(path.Join("api", "v1", api))
	if err != nil {
	   return "", err
	}
	
	// append token to the URL if the users uses token based authentication
	secret, err := apiSecret()
	if err != nil {
		return "", err
	}
	hasToken := strings.HasPrefix(secret, "token=")
	if hasToken { // users uses token based authentication
		token := secret[6:len(secret)] // drop the token= prefix
		match, _ := regexp.MatchString("^[a-z_0-9]{0,10}-[a-f0-9]{16}$", token) // subjectName is up to ten lowercase letters, numbers or underscore, followed by '-' and followed by a 16 hex shasum digest
		if match {
			q := u.Query() // Get a copy of the query values.
			q.Add("token", token)
			u.RawQuery = q.Encode() // Encode and assign back to the original query.
		} else {
			return "", fmt.Errorf("Not a valid Nightscout token: %s. Expected ^token=[a-z_0-9]{0,10}-[a-f0-9]{16}$ as Nightscout secret", secret)
		}
	}
	return siteURL.ResolveReference(u).String(), nil
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
	secret, err := apiSecret()
	if err != nil {
		return err
	}
	req.Header.Add("api-secret", secret)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	return nil
}

// Get performs a GET operation on a Nightscout API.
func Get(api string, result interface{}) error {
	data, err := RestOperation("GET", api, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

// Upload performs a POST or PUT operation on a Nightscout API.
func Upload(op string, api string, param interface{}) error {
	switch op {
	case "POST", "PUT":
		if param == nil {
			log.Panicf("%s %s operation with no data", op, api)
		}
	default:
		log.Panicf("%s %s used for upload", op, api)
	}
	_, err := RestOperation(op, api, param)
	return err
}

func indentJSON(data []byte) string {
	buf := bytes.Buffer{}
	err := json.Indent(&buf, data, "", "  ")
	if err != nil {
		log.Printf("json.Indent error: %v", err)
		return string(data)
	}
	return buf.String()
}

// Hostname returns the host name.
func Hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

// Username returns the user name.
func Username() string {
	u := os.Getenv("USER")
	if len(u) == 0 {
		return "unknown"
	}
	return u
}

// Device returns the Nightscout device name.
func Device() string {
	u := os.Getenv(deviceEnvVar)
	if len(u) == 0 {
		return "openaps://" + Hostname()
	}
	return u
}

// JSON marshals the given data in indented form
// and returns it as a string.
func JSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return indentJSON(data)
}
