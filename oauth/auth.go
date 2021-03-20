package oauth

import (
	"miikka.xyz/rs/oauth/github"
	"miikka.xyz/rs/store"
)

// GetGithub adds support for GitHub authentication
func GetGithub(db *store.DB) *github.Auth {
	github := &github.Auth{DB: db}
	return github
}
