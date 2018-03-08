package module

import "strings"

// Locator is a string specifying a particular dependency and revision
type Locator string

// Dependency represents a code library brought in by running a Build
type Dependency interface {
	Locator() Locator
}

func normalizeGitURL(project string) string {
	// Remove fetcher prefix (in case project is derived from splitting a locator on '$')
	noFetcherPrefix := strings.TrimPrefix(project, "git+")

	// Normalize Git URL format
	noGitExtension := strings.TrimSuffix(noFetcherPrefix, ".git")
	handleGitHubSSH := strings.Replace(noGitExtension, "git@github.com:", "github.com/", 1)

	// Remove protocols
	noHTTPPrefix := strings.TrimPrefix(handleGitHubSSH, "http://")
	noHTTPSPrefix := strings.TrimPrefix(noHTTPPrefix, "https://")

	return noHTTPSPrefix
}

// MakeLocator creates a locator string given a package and revision
func MakeLocator(fetcher string, project string, revision string) string {
	if fetcher != "git" {
		return fetcher + "+" + project + "$" + revision
	}
	return "git+" + normalizeGitURL(project) + "$" + revision
}
