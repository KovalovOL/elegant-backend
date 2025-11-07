package user

type UserGoogleResp struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type CreateUser struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	GitHubUrl  string `json:"github_url"`
	LikedinUrl string `json:"linkedin_url"`
	Bio        string `json:"bio"`
}

type User struct {
	CreateUser
	ID string `json:"id"`
}
