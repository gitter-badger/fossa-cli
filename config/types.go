package config

import (
	"strings"
	"time"

	logging "github.com/op/go-logging"

	"github.com/fossas/fossa-cli/module"
)

var configLogger = logging.MustGetLogger("config")

// DefaultConfig specifies the config for the default command
type DefaultConfig struct {
	Build bool
}

// AnalyzeConfig specifies the config for the analyze command
type AnalyzeConfig struct {
	Output          bool
	AllowUnresolved bool
}

// BuildConfig specifies the config for the build command
type BuildConfig struct {
	Force bool
}

// TestConfig specifies the config for the test command
type TestConfig struct {
	Timeout time.Duration
}

// UploadConfig specifies the config for the upload command
type UploadConfig struct {
	Locators bool
	Data     string
}

// ReportConfig specifies the config for the report command
type ReportConfig struct {
	Type string // Either "dependencies" or "licenses"
}

// CLIConfig specifies the config available to the cli
type CLIConfig struct {
	APIKey   string
	Fetcher  string
	Project  string
	Revision string
	Endpoint string
	Modules  []module.Config
	Debug    bool

	DefaultCmd DefaultConfig
	AnalyzeCmd AnalyzeConfig
	BuildCmd   BuildConfig
	TestCmd    TestConfig
	UploadCmd  UploadConfig
	ReportCmd  ReportConfig

	ConfigFilePath string
}

// MakeLocator creates a locator string given a package and revision
func MakeLocator(fetcher string, project string, revision string) string {
	if fetcher != "git" {
		return fetcher + "+" + project + "$" + revision
	}

	return "git+" + normalizePackageSpec(project) + "$" + revision
}

func normalizePackageSpec(project string) string {
	// Remove fetcher prefix (in case project is derived from splitting a locator on '$')
	noFetcherPrefix := strings.TrimPrefix(project, "git+")

	// Normalise Git URL format
	noGitExtension := strings.TrimSuffix(noFetcherPrefix, ".git")
	handleGitHubSSH := strings.Replace(noGitExtension, "git@github.com:", "github.com/", 1)

	// Remove protocols
	noHTTPPrefix := strings.TrimPrefix(handleGitHubSSH, "http://")
	noHTTPSPrefix := strings.TrimPrefix(noHTTPPrefix, "https://")

	return noHTTPSPrefix
}
