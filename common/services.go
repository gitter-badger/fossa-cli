package common

import (
	"github.com/briandowns/spinner"
)

// Services is a container for all services. Each service is a self-contained
// interface of side-effecting functions. No service may rely on any other
// service.
type Services interface {
	LoggingService() LoggingService
	APIService() APIService
	ExecService() ExecService
	FileService() FileService
}

// A LoggingService implements all logging functionality. All output is printed
// to STDERR, except for the specific case of `Printf()`, which prints to
// STDOUT. STDOUT should be reserved only for actual meaningful output, and
// STDERR should be used for all other user-interactive output.
type LoggingService interface {
	// Returns the spinner. There should only be at most one spinner active at
	// one time.
	Spinner() *spinner.Spinner

	// Debug statements are meant for diagnosing and resolving issues. They are
	// not shown unless `--debug` is specified. When `--debug` is set, all
	// statements are logged with extra debug information.
	Debug()
	Debugf()

	// Notices inform the user of a non-error condition that is important.
	Notice()
	Noticef()

	// Warnings inform the user of a non-fatal error condition.
	Warning()
	Warningf()

	// Fatals cause the program to exit with a non-zero exit code. They inform
	// the user of fatal error conditions.
	Fatal()
	Fatalf()

	// Printing sends output to STDOUT.
	Printf()
}

// ErrTimeout is an error caused by a connection timeout.
type ErrTimeout = error

// An APIService implements low-level HTTP functionality. The API package
// implements a high-level interface on top of this.
type APIService interface {
	// Initialize configures the APIService with a default server and API key.
	Initialize(server, APIKey string) error

	// These functions are for convenience; for most cases, the server and API
	// key don't change over the life of the command.
	Get(URL string, body []byte) (res string, statusCode int, err error)
	Post(URL string, body []byte) (res string, statusCode int, err error)

	GetJSON(URL string, body []byte, v interface{}) (statusCode int, err error)
	PostJSON(URL string, body []byte, v interface{}) (statusCode int, err error)

	// This is the underlying implementation of the APIService's functionality.
	MakeAPIRequest(server, method, URL, APIKey string, body []byte) (res []byte, statusCode int, err error)
}

// A WhichResolver is a function which determines the version and existence of
// an external binary given a command.
type WhichResolver func(cmd string) (version string, ok bool)

// An ExecService implements calls to external commands.
type ExecService interface {
	RunCWD(cmd string) (stdout string, stderr string, err error)
	Run(dir string, cmd string) (stdout string, stderr string, err error)

	Which(args string, candidates ...string) (cmd string, version string, err error)
	WhichWithResolver(args string, resolver WhichResolver, candidates ...string) (cmd string, version string, err error)
}

// An UnmarshalFunc unmarshals a particular data format. `json.Unmarshal` is an
// example of this.
type UnmarshalFunc func(data []byte, v interface{}) error

// A FileService implements filesystem interaction.
type FileService interface {
	HasFile(file string) bool
	HasFolder(file string) bool
	ReadFile(file string) ([]byte, error)
	ReadJSON(file string, v interface{}) error
	ReadUnmarshal(file string, unmarshal UnmarshalFunc, v interface{}) error
}
