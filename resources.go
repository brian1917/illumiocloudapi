package illumiocloudapi

// Resource represents a cloud resource with its associated metadata.
type Resource struct {
	AccountID   string `json:"account_id"`
	Cloud       string `json:"cloud"`
	ObjectType  string `json:"object_type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	AccountName string `json:"account_name"`
}

// ResourcesRespBody represents the response body for a request fetching cloud resources.
type ResourcesRespBody struct {
	CloudResources []Resource `json:"items"`
	NextPageToken  string     `json:"next_page_token"`
	Page           int        `json:"page"`
}

// ResourcesReqBody represents the request body for a request fetching cloud resources.
type ResourcesReqBody struct {
	PageToken string `json:"page_token"`
}

func (t *Tenant) GetResources(queryParameters map[string]string) (apiResponses []ApiResponse, err error) {

	var resourcesRespBody ResourcesRespBody
	var resources []Resource

	// Start by making the call without any body
	api, err := t.Post("inventory/resources", nil, &resourcesRespBody)
	if err != nil {
		return apiResponses, err
	}
	apiResponses = append(apiResponses, api)
	resources = append(resources, resourcesRespBody.CloudResources...)

	// Iterate while we have more pages
	for resourcesRespBody.NextPageToken != "" {
		api, err = t.Post("inventory/resources", nil, &resourcesRespBody)
		if err != nil {
			return apiResponses, err
		}
		apiResponses = append(apiResponses, api)
		resources = append(resources, resourcesRespBody.CloudResources...)
	}

	t.Resources = resources

	return apiResponses, err
}
