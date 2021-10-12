package nightscout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Website struct {
	URL      *url.URL
	Client   *http.Client
	Token    string
	noUpload bool
	verbose  bool
}

const (
	siteEnvVar      = "NIGHTSCOUT_SITE"
	apiSecretEnvVar = "NIGHTSCOUT_API_SECRET"
	deviceEnvVar    = "NIGHTSCOUT_DEVICE"
)

func Site(site string) (*Website, error) {
	u, err := url.Parse(site)
	if err != nil {
		return nil, err
	}
	return &Website{
		URL:    u,
		Client: &http.Client{},
	}, nil
}

func DefaultSite() (*Website, error) {
	site := os.Getenv(siteEnvVar)
	if len(site) == 0 {
		return nil, fmt.Errorf("%s is not set", siteEnvVar)
	}
	return Site(site)
}

func (w *Website) String() string {
	return w.URL.String()
}

func (w *Website) APISecret() (string, error) {
	secret := os.Getenv(apiSecretEnvVar)
	if len(secret) == 0 {
		if len(w.Token) == 0 {
			return "", fmt.Errorf("%s is not set", apiSecretEnvVar)
		} else {
			secret = w.Token
		}
	}
	return secret, nil
}

// Verbose returns the value of the verbose flag.
func (w *Website) Verbose() bool {
	return w.verbose
}

// SetVerbose sets the value of the verbose flag.
func (w *Website) SetVerbose(flag bool) {
	w.verbose = flag
}

// NoUpload returns the value of the noUpload flag.
func (w *Website) NoUpload() bool {
	return w.noUpload
}

// SetNoUpload sets the value of the noUpload flag.
func (w *Website) SetNoUpload(flag bool) {
	w.noUpload = flag
}

func (w *Website) restOperation(op string, api string, data interface{}, result interface{}) error {
	switch op {
	case "GET":
		if data != nil {
			log.Panicf("GET %s operation with data", api)
		}
	case "POST", "PUT":
		if data == nil {
			log.Panicf("%s %s operation with no data", op, api)
		}
	default:
		log.Panicf("unsupported %s %s operation", op, api)
	}
	req, err := w.makeRequest(op, api, data)
	if err != nil {
		return err
	}
	if w.verbose || w.noUpload {
		u := req.URL.String()
		q, err := url.QueryUnescape(u)
		if err != nil {
			q = u
		}
		log.Printf("%s %s", op, q)
		if data != nil {
			log.Print(JSON(data))
		}
	}
	if w.noUpload && op != "GET" {
		return nil
	}
	resp, err := w.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	if code != http.StatusOK {
		return fmt.Errorf("%v: %d %s", req.URL, code, http.StatusText(code))
	}
	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
	}
	if w.verbose && err == nil && result != nil {
		log.Print(JSON(result))
	}
	return err
}

func (w *Website) makeRequest(op string, api string, data interface{}) (*http.Request, error) {
	u, err := w.makeURL(op, api)
	if err != nil {
		return nil, err
	}
	r, err := makeReader(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(op, u, r)
	if err != nil {
		return nil, err
	}
	err = w.addHeaders(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (w *Website) makeURL(op string, api string) (string, error) {
	u, err := url.Parse(api)
	if err != nil {
		return "", err
	}
	secret, err := w.APISecret()
	if err != nil {
		return "", err
	}
	if usesTokenAuth(secret) {
		// Validate token.
		token := secret[len("token="):]
		if !validToken.MatchString(token) {
			return "", fmt.Errorf("invalid Nightscout token %q", secret)
		}
		// Append token to the URL parameters.
		q := u.Query()
		q.Add("token", token)
		u.RawQuery = q.Encode()
	}
	return w.URL.ResolveReference(u).String(), nil
}

// Auth token must be of the form <subject name>-<hash code>,
// where the subject name is up to ten lowercase letters, digits, or underscores
// and the hash code is the first 16 hex digits of the SHA-1 digest of the API secret plus Mongo ObjectID.
var validToken = regexp.MustCompile(`^[a-z_0-9]{0,10}-[a-f0-9]{16}$`)

func usesTokenAuth(secret string) bool {
	return strings.HasPrefix(secret, "token=")
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

func (w *Website) addHeaders(req *http.Request) error {
	secret, err := w.APISecret()
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	if !usesTokenAuth(secret) {
		req.Header.Add("api-secret", secret)
	}
	return nil
}

// Get performs a GET operation on a Nightscout API.
func (w *Website) Get(api string, result interface{}) error {
	return w.restOperation("GET", api, nil, result)
}

// Upload performs a POST operation on a Nightscout API.
func (w *Website) Upload(api string, data interface{}) error {
	return w.restOperation("POST", api, data, nil)
}

// Put performs a PUT operation on a Nightscout API.
func (w *Website) Put(api string, data interface{}) error {
	return w.restOperation("PUT", api, data, nil)
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
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}
