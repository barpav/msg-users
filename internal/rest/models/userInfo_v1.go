package models

// Schema: userInfo.v1
type UserInfoV1 struct {
	Name    string `json:"name"`
	Picture string `json:"picture,omitempty"`
}
