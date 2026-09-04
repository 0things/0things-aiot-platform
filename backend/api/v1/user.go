package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
} //@name ApiRegisterRequest

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
} //@name ApiLoginRequest
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
} //@name ApiLoginResponseData

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
} //@name ApiUpdateProfileRequest
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" example:"alan@gmail.com"`
} //@name ApiGetProfileResponseData
