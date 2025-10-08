package illumiocloudapi

import (
	"fmt"
	"time"

	"github.com/brian1917/illumioapi/v2"
	"github.com/brian1917/workloader/utils"
)

// Resource represents a cloud resource with its associated metadata.
type Resource struct {
	AccountID   string      `json:"account_id"`
	Cloud       string      `json:"cloud"`
	ObjectType  string      `json:"object_type"`
	Category    string      `json:"category"`
	Subcategory string      `json:"subcategory"`
	AccountName string      `json:"account_name"`
	Name        string      `json:"name"`
	Id          string      `json:"id"`
	Relations   []Relations `json:"relations,omitempty"`
}

// ResourcesPostRequest represents the request body for a request fetching cloud resources.
type ResourcesPostRequest struct {
	PageToken  string   `json:"page_token"`
	ObjectType []string `json:"object_type,omitempty"`
}

// ResourcesPostResponse represents the response body for a request fetching cloud resources.
type ResourcesPostResponse struct {
	CloudResources []Resource `json:"items"`
	NextPageToken  string     `json:"next_page_token"`
	Page           int        `json:"page"`
	TotalSize      int        `json:"total_size"`
}

type LabelingPostRequest struct {
	LabelAssignments []LabelAssignment `json:"label_assignments"`
}

type LabelingPostResponse struct {
	FailedResources []FailedResources `json:"failed_resources"`
}
type FailedResources struct {
	CspID string `json:"csp_id"`
	Error string `json:"error"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LabelAssignment struct {
	CspID  string   `json:"csp_id"`
	Add    *[]Label `json:"add,omitempty"`
	Remove *[]Label `json:"remove,omitempty"`
}

type Relations struct {
	ID          string     `json:"id"`
	CspID       string     `json:"csp_id"`
	AccountID   string     `json:"account_id"`
	TenantID    string     `json:"tenant_id"`
	Cloud       string     `json:"cloud"`
	Name        string     `json:"name"`
	ObjectType  string     `json:"object_type"`
	Category    string     `json:"category"`
	Region      string     `json:"region"`
	State       string     `json:"state"`
	CreatedAt   time.Time  `json:"created_at"`
	Properties  Properties `json:"properties,omitempty"`
	Subcategory string     `json:"subcategory"`
	AccountName string     `json:"account_name"`
}

type Properties struct {
	AcceptorCspID  string `json:"acceptor_csp_id"`
	RequesterCspID string `json:"requester_csp_id"`
	ResourceGroup  string `json:"resource_group"`
}

func (t *Tenant) GetResources(requestBody ResourcesPostRequest) (apiResponses map[string]illumioapi.APIResponse, err error) {

	var resourcesRespBody ResourcesPostResponse
	var resources []Resource
	apiResponses = make(map[string]illumioapi.APIResponse)

	// Start by making the call without any body
	api, err := t.Post("inventory/resources", requestBody, &resourcesRespBody)
	if err != nil {
		return apiResponses, err
	}
	apiResponses[fmt.Sprintf("inventory/resources-page%d", 1)] = api
	utils.LogInfof(true, "%d cloud resources", resourcesRespBody.TotalSize)
	resources = append(resources, resourcesRespBody.CloudResources...)
	utils.LogInfof(true, "%d resources downloaded from page %d", len(resourcesRespBody.CloudResources), resourcesRespBody.Page)

	// Iterate while we have more pages
	for resourcesRespBody.NextPageToken != "" {
		requestBody.PageToken = resourcesRespBody.NextPageToken
		resourcesRespBody = ResourcesPostResponse{}
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

func (t *Tenant) LabelResources(labelAssignments []LabelAssignment) (apiResponse illumioapi.APIResponse, err error) {
	labelingPostRequest := LabelingPostRequest{
		LabelAssignments: labelAssignments,
	}
	labelingPostResponse := LabelingPostResponse{}
	apiResponse, err = t.Post("label_assignments", labelingPostRequest, &labelingPostResponse)
	if len(labelingPostResponse.FailedResources) > 0 {
		return apiResponse, fmt.Errorf("failed to label %d resources", len(labelingPostResponse.FailedResources))
	}
	if err != nil {
		return apiResponse, err
	}
	return apiResponse, nil
}
