package illumiocloudapi

import "strings"

// GetAzureSubscription extracts and returns the subscription ID from a given resource ID.
func GetAzureSubscription(resourceId string) string {
	x := strings.Split(resourceId, "/")
	for i, v := range x {
		if v == "subscriptions" && len(x) > i+1 {
			return x[i+1]
		}
	}
	return ""
}

// GetAzureResourceGroup extracts and returns the resource Group from a given resource ID.
func GetAzureResourceGroup(resourceId string) string {
	x := strings.Split(resourceId, "/")
	for i, v := range x {
		if v == "resourceGroups" && len(x) > i+1 {
			return x[i+1]
		}
	}
	return ""
}

// GetAzureResourceName extracts the resource name from a given resource ID.
func GetAzureResourceName(resourceId string) string {
	x := strings.Split(resourceId, "/")
	if len(x) == 0 {
		return ""
	}
	return x[len(x)-1]
}
