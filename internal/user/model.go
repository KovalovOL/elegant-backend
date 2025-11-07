package user

type UserGoogleResp struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type CreaeteUser struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	GitHubUrl  string `json:"github_url"`
	LikedinUrl string `json:"linkedin_url"`
	Bio        string `json:"bio"`
}

type User struct {
	CreaeteUser
	ID string `json:"id"`
}
