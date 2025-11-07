package cv


type CreateCVInput struct {
	Title	   string `json:"title"`
	Position    string `json:"position"`
	Summary     string `json:"summary"`
	Skills      string `json:"skills"`
	Experience  string `json:"experience"`
	Education   string `json:"education"`
}

type CreateCV struct {
	CreateCVInput
	UserID      int    `json:"user_id"`
}

type CV struct {
	CreateCV
	CVID int `json:"cv_id"`
}
