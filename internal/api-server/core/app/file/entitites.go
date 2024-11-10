package file

// ImageConfig represents the configuration for image uploads.
type ImageConfig struct {
	NumberLimits    int                       `json:"number_limits"`
	TransferMethods []FileTransferMethod      `json:"transfer_methods"`
	Detail          *ImagePromptMessageDetail `json:"detail"`
}

// FileExtraConfig represents additional configuration for file uploads.
type FileExtraConfig struct {
	ImageConfig          *ImageConfig         `json:"image_config"`
	AllowedFileTypes     []FileType           `json:"allowed_file_types"`
	AllowedExtensions    []string             `json:"allowed_extensions"`
	AllowedUploadMethods []FileTransferMethod `json:"allowed_upload_methods"`
	NumberLimits         int                  `json:"number_limits"`
}

// File represents the file entity used in the system.
type File struct {
	DifyModelIdentity string             `json:"dify_model_identity"`
	ID                *string            `json:"id,omitempty"` // Message file ID, optional
	TenantID          string             `json:"tenant_id"`
	Type              FileType           `json:"type"`
	TransferMethod    FileTransferMethod `json:"transfer_method"`
	RemoteURL         *string            `json:"remote_url,omitempty"`
	RelatedID         *string            `json:"related_id,omitempty"`
	Filename          *string            `json:"filename,omitempty"`
	Extension         *string            `json:"extension,omitempty"` // File extension, should contain a dot
	MimeType          *string            `json:"mime_type,omitempty"`
	Size              int                `json:"size"`
	ExtraConfig       *FileExtraConfig   `json:"extra_config,omitempty"`
}

// ImagePromptMessageDetail is a placeholder struct for the "Detail" field in ImageConfig.
type ImagePromptMessageDetail struct {
	// Define fields as needed.
}
