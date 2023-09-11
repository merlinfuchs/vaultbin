package wire

import "time"

type InternalPasteCreateRequest struct {
	Content  string `form:"content"`
	Language string `form:"language"`
}

type InternalPasteDuplicateRequest struct {
	Content string `form:"content"`
}

type InternalPasteRawRequest struct {
	ID string `form:"id"`
}

type APIPasteCreateRequest struct {
	Content    string `form:"content" json:"content"`
	Language   string `form:"language" json:"language"`
	Expiration int    `from:"expiration" json:"expiration"`
}

type APIPaste struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	ViewCount int       `json:"view_count"`
}

type APIPasteGetResponseData APIPaste

type APIPasteGetResponse APIResponse[APIPasteGetResponseData]

type APIPasteCreateResponseData APIPaste

type APIPasteCreateResponse APIResponse[APIPasteCreateResponseData]
