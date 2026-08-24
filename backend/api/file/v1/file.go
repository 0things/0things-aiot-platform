// Package filev1 owns HTTP contracts for file uploads.
package filev1

// UploadOTAFileResponse contains the immutable metadata used to create an OTA
// package after its binary has been stored in object storage.
type UploadOTAFileResponse struct {
	FileURL  string `json:"fileUrl"`
	FileSize int64  `json:"fileSize"`
	Checksum string `json:"checksum"`
} //@name FileUploadOTAFileResponse
