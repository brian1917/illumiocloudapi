package illumiocloudapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Post sends a POST request to the cloud tenant
func (t *Tenant) Post(endpoint string, object, createdObject interface{}) (api ApiResponse, err error) {
	// Build the API URL
	apiURL, err := url.Parse(fmt.Sprintf("https://%s/api/v1/%s", t.Url, endpoint))
	if err != nil {
		return api, err
	}

	// Create payload
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return api, err
	}

	// Call the API
	api, err = t.HttpReq("POST", apiURL.String(), jsonBytes)
	api.ReqBody = string(jsonBytes)
	if err != nil {
		return api, err
	}

	// Unmarshal new label
	json.Unmarshal([]byte(api.RespBody), &createdObject)

	return api, nil
}
