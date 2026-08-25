package v1

type OrganizationItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	IsCurrent bool   `json:"is_current"`
} //@name ApiOrganizationItem

type SwitchOrgRequest struct {
	OrgId int64 `json:"org_id" binding:"required"`
} //@name ApiSwitchOrgRequest

type SwitchOrgResponseData struct {
	AccessToken string `json:"accessToken"`
} //@name ApiSwitchOrgResponseData
