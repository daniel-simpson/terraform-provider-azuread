package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Client contains all configuration and objects necessary to interact with the Microsoft Graph API
type Client struct {
	baseURL      *url.URL
	tenantID     string
	clientID     string
	clientSecret string
	bearerToken  string
	tokenExpires *time.Time
	httpClient   *http.Client
}

type bearerTokenRequest struct {
	scope        string `json:"scope"`
	grantType    string `json:"grant_type"`
	clientID     string `json:"client_id"`
	clientSecret string `json:"client_secret"`
}

type bearerTokenResponse struct {
	tokenType    string `json:"token_type"`
	expiresIn    int    `json:"expires_in"`
	extExpiresIn int    `json:"ext_expires_in"`
	accessToken  string `json:"access_token"`
}

func (c *Client) createClient(baseURL string) error {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Timeout: time.Minute * 2,
	}

	c.httpClient = httpClient
	c.baseURL = parsedBaseURL

	return nil
}

// newRequest scaffolds a bare http.Request with the necessary bearer token
func (c *Client) newRequest(method string, relativeURL string, body io.Reader) (*http.Request, error) {
	rel, parseError := url.Parse(relativeURL)
	if parseError != nil {
		return nil, parseError
	}

	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.baseURL.ResolveReference(rel)

	req, requestError := http.NewRequest(method, u.String(), nil)
	if requestError != nil {
		return nil, requestError
	}

	bearerToken, tokenError := c.getBearerToken()
	if tokenError != nil {
		return nil, tokenError
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken)

	return req, nil
}

func (c *Client) getBearerToken() (string, error) {
	if c.bearerToken == "" || c.tokenExpires.IsZero() || c.tokenExpires.Before(time.Now()) {
		payloadBuffer, payloadBufferError := getPayloadBuffer(bearerTokenRequest{
			scope:        "https://graph.microsoft.com/.default",
			grantType:    "client_credentials",
			clientID:     c.clientID,
			clientSecret: c.clientSecret,
		})

		if payloadBufferError != nil {
			return "", payloadBufferError
		}

		req, requestError := http.NewRequest("POST", fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", c.tenantID), payloadBuffer)

		if requestError != nil {
			return "", requestError
		}

		resp, responseError := c.httpClient.Do(req)

		if responseError != nil {
			return "", responseError
		}

		defer resp.Body.Close()

		body := bearerTokenResponse{}

		decodeError := parseResponse(resp, &body)
		if decodeError != nil {
			return "", decodeError
		}

		expiry := time.Now().Add(time.Duration(body.expiresIn) * time.Second)

		c.tokenExpires = &expiry
		c.bearerToken = body.accessToken
	}

	return fmt.Sprintf("Bearer %s", c.bearerToken), nil
}

func getPayloadBuffer(payload interface{}) (*bytes.Buffer, error) {
	payloadStr, marshalErrorson := json.Marshal(payload)
	if marshalErrorson != nil {
		return nil, marshalErrorson
	}

	return bytes.NewBuffer(payloadStr), nil
}

func parseResponse(httpResp *http.Response, body interface{}) error {

	var bodyBytes []byte
	bodyBytes, responseError := ioutil.ReadAll(httpResp.Body)

	if responseError != nil {
		return errors.Wrap(responseError, "parseResponse failed to read response")
	}

	decodeError := json.Unmarshal(bodyBytes, &body)

	if decodeError != nil {
		return errors.Wrapf(decodeError, "parseResponse failed for input: %s", string(bodyBytes))
	}

	if httpResp.StatusCode < 200 || 299 < httpResp.StatusCode {
		return errors.New(fmt.Sprintf("Invalid response code %v, unable to retrieve subnet.  Body: %v", httpResp.StatusCode, string(bodyBytes)))
	}

	return nil
}
