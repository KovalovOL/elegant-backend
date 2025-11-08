package tag

type CreateTag struct {
	Name string `json:"name"`
}

type Tag struct {
	CreateTag
	TagID int `json:"tag_id"`
}
