package models

// Client defines cleint object
type Client struct {
	ID         string `json:"_id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Note       string `json:"note"`
	IsArchived bool   `json:"is_archived"`
}
