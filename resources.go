package illumiocloudapi

import (
	"fmt"

	"github.com/brian1917/illumioapi/v2"
	"github.com/brian1917/workloader/utils"
)

// Resource represents a cloud resource with its associated metadata.
type Resource struct {
	AccountID   string `json:"account_id"`
	Cloud       string `json:"cloud"`
	ObjectType  string `json:"object_type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	AccountName string `json:"account_name"`
	Name        string `json:"name"`
	Id          string `json:"id"`
}

// ResourcesRespBody represents the response body for a request fetching cloud resources.
type ResourcesRespBody struct {
	CloudResources []Resource `json:"items"`
	NextPageToken  string     `json:"next_page_token"`
	Page           int        `json:"page"`
	TotalSize      int        `json:"total_size"`
}

// ResourcesReqBody represents the request body for a request fetching cloud resources.
type ResourcesReqBody struct {
	PageToken string `json:"page_token"`
}

func (t *Tenant) GetResources(queryParameters map[string]string) (apiResponses map[string]illumioapi.APIResponse, err error) {

	var resourcesRespBody ResourcesRespBody
	var resources []Resource
	apiResponses = make(map[string]illumioapi.APIResponse)

	// Start by making the call without any body
	api, err := t.Post("inventory/resources", ResourcesReqBody{}, &resourcesRespBody)
	if err != nil {
		return apiResponses, err
	}
	apiResponses[fmt.Sprintf("inventory/resources-page%d", 1)] = api
	utils.LogInfof(true, "%d cloud resources", resourcesRespBody.TotalSize)
	resources = append(resources, resourcesRespBody.CloudResources...)
	utils.LogInfof(true, "%d resources downloaded from page %d", len(resourcesRespBody.CloudResources), resourcesRespBody.Page)

	// Iterate while we have more pages
	for resourcesRespBody.NextPageToken != "" {
		requestBody := ResourcesReqBody{PageToken: resourcesRespBody.NextPageToken}
		resourcesRespBody = ResourcesRespBody{}
		api, err = t.Post("inventory/resources", requestBody, &resourcesRespBody)
		if err != nil {
			return apiResponses, err
		}
		apiResponses[fmt.Sprintf("inventory/resources-page%d", resourcesRespBody.Page)] = api
		resources = append(resources, resourcesRespBody.CloudResources...)
		utils.LogInfof(true, "%d resources downloaded from page %d", len(resourcesRespBody.CloudResources), resourcesRespBody.Page)
	}

	t.Resources = resources

	return apiResponses, err
}
