// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

// FileType represents different file types in the system.
type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeDocument FileType = "document"
	FileTypeAudio    FileType = "audio"
	FileTypeVideo    FileType = "video"
	FileTypeCustom   FileType = "custom"
)

// FileTransferMethod represents the file transfer methods available.
type FileTransferMethod string

const (
	FileTransferMethodRemoteURL FileTransferMethod = "remote_url"
	FileTransferMethodLocalFile FileTransferMethod = "local_file"
	FileTransferMethodToolFile  FileTransferMethod = "tool_file"
)

// FileBelongsTo represents who the file belongs to.
type FileBelongsTo string

const (
	FileBelongsToUser      FileBelongsTo = "user"
	FileBelongsToAssistant FileBelongsTo = "assistant"
)

// FileAttribute represents various file attributes.
type FileAttribute string

const (
	FileAttributeType           FileAttribute = "type"
	FileAttributeSize           FileAttribute = "size"
	FileAttributeName           FileAttribute = "name"
	FileAttributeMimeType       FileAttribute = "mime_type"
	FileAttributeTransferMethod FileAttribute = "transfer_method"
	FileAttributeURL            FileAttribute = "url"
	FileAttributeExtension      FileAttribute = "extension"
)

// ArrayFileAttribute represents attributes for arrays of files.
type ArrayFileAttribute string

const (
	ArrayFileAttributeLength ArrayFileAttribute = "length"
)
