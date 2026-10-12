package v1

type AuthMeResp struct {
	Subject        string `json:"sub"`
	OrganizationID string `json:"organization_id"`
}
