package illumiocloudapi

import (
	"fmt"

	"github.com/brian1917/illumioapi/v2"
)

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LabelAssignment struct {
	CspID  string   `json:"csp_id"`
	Add    *[]Label `json:"add,omitempty"`
	Remove *[]Label `json:"remove,omitempty"`
}

type FailedResources struct {
	CspID string `json:"csp_id"`
	Error string `json:"error"`
}

type LabelingPostRequest struct {
	LabelAssignments []LabelAssignment `json:"label_assignments"`
}

type LabelingPostResponse struct {
	FailedResources []FailedResources `json:"failed_resources"`
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

// Returns the value of a label by its key. If the key does not exist, it returns an empty string.
func (r *Resource) GetLabelValueByKey(key string) (value string) {
	for _, label := range r.Labels {
		if label.Key == key {
			return label.Value
		}
	}
	return ""
}
