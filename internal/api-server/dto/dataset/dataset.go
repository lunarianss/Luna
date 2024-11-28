// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

// FileUploadConfigurationResponse represents the configuration for file uploads.
type FileUploadConfigurationResponse struct {
	FileSizeLimit           int `json:"file_size_limit"`            // The maximum size limit for files (in MB).
	BatchCountLimit         int `json:"batch_count_limit"`          // The maximum number of files in a batch.
	ImageFileSizeLimit      int `json:"image_file_size_limit"`      // The maximum size limit for image files (in MB).
	VideoFileSizeLimit      int `json:"video_file_size_limit"`      // The maximum size limit for video files (in MB).
	AudioFileSizeLimit      int `json:"audio_file_size_limit"`      // The maximum size limit for audio files (in MB).
	WorkflowFileUploadLimit int `json:"workflow_file_upload_limit"` // The maximum file upload limit for workflows.
}

func NewFileUploadConfigurationResponse() *FileUploadConfigurationResponse {
	return &FileUploadConfigurationResponse{
		FileSizeLimit:           15,
		BatchCountLimit:         5,
		ImageFileSizeLimit:      10,
		VideoFileSizeLimit:      100,
		AudioFileSizeLimit:      50,
		WorkflowFileUploadLimit: 10,
	}
}
