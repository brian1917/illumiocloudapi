package illumiocloudapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/brian1917/illumioapi/v2"
)

// A tenant represents an Illumio cloud tenant
type Tenant struct {
	TenantID  string     `json:"tenant_id"`
	Cookie    string     `json:"cookie"`
	Resources []Resource `json:"resources,omitempty"`
	ClientID  string     `json:"client_id,omitempty"`
	Key       string     `json:"key,omitempty"`
}

// HttpReq makes an API call to an Illumio Cloud tenant with sepcified options
// action must be GET, POST, PUT, or DELETE.
// url is the full endpoint being called.
// PUT and POST methods should have a body that is JSON run through the json.marshal function so it's a []byte.
func (t *Tenant) HttpReq(action, url string, body []byte) (illumioapi.APIResponse, error) {

	// Setup the http client
	client := &http.Client{}
	req, err := http.NewRequest(action, url, bytes.NewBuffer(body))
	if err != nil {
		return illumioapi.APIResponse{}, err
	}

	// Set headers and authenticate
	req.Header.Add("X-Tenant-Id", t.TenantID)
	req.Header.Add("Content-Type", "application/json")
	if t.Cookie == "" && (t.ClientID == "" || t.Key == "") {
		return illumioapi.APIResponse{}, fmt.Errorf("either cookie or client_id and key must be set for authentication")
	}
	if t.ClientID != "" && t.Key != "" {
		req.SetBasicAuth(t.ClientID, t.Key)
	} else if t.Cookie != "" {
		req.Header.Add("Cookie", t.Cookie)
	}

	// Make the http request
	response, err := client.Do(req)
	if err != nil {
		return illumioapi.APIResponse{}, err
	}
	defer response.Body.Close()

	// Read the response
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return illumioapi.APIResponse{}, err
	}

	apiResp := illumioapi.APIResponse{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Request:    req,
		RespBody:   string(responseData),
	}
	return apiResp, nil
}
