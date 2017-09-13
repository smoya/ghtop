package contributor

import "github.com/google/go-github/github"

const (
	// VCSGithub is the name given to github as VCS.
	VCSGithub = "github"
)

// Contributor represents a user from a vcs system.
type Contributor struct {
	Username   string `json:"username"`
	VCS        string `json:"vcs"`
	ProfileURL string `json:"profile_url"`
	AvatarURL  string `json:"avatar_url"`
}

// FromGithub adapts a Github user to a Contributor.
func FromGithub(u github.User) Contributor {
	return Contributor{
		Username:   u.GetLogin(),
		VCS:        VCSGithub,
		ProfileURL: u.GetURL(),
		AvatarURL:  u.GetAvatarURL(),
	}
}
