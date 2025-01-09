package dto

type PreviewFileQuery struct {
	Timestamp    int64  `json:"timestamp" form:"timestamp" validate:"required"`
	Nonce        string `json:"nonce" form:"nonce" validate:"required"`
	Sign         string `json:"sign" form:"sign" validate:"required"`
	AsAttachment bool   `json:"as_attachment" form:"as_attachment"`
}

type PreviewFileUri struct {
	Filename string `uri:"filename" json:"filename" validate:"required"`
}
