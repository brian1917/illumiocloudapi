package illumiocloudapi

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/brian1917/illumioapi/v2"
)

// Set the constant cloud base url
const CloudBaseFqdn = "cloud.illum.io"

// Post sends a POST request to the cloud tenant
func (t *Tenant) Post(endpoint string, object, createdObject interface{}) (api illumioapi.APIResponse, err error) {
	// Build the API URL
	apiURL, err := url.Parse(fmt.Sprintf("https://%s/api/v1/%s", CloudBaseFqdn, endpoint))
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

	// Unmarshal response
	json.Unmarshal([]byte(api.RespBody), &createdObject)

	return api, nil
}
