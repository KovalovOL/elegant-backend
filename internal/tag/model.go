package tag

type TagCreate struct {
	Name  string `json:"name"`
}

type Tag struct {
	TagCreate
	TagID int    `json:"tag_id"`
}

