package module

import (
	"os"
	"path/filepath"
	"strings"

	logging "github.com/op/go-logging"
)

var moduleLogger = logging.MustGetLogger("module")

// A Module is a unit of buildable code within a project.
type Module struct {
	Name string
	Type Type
	// Target is the entry point or manifest path for the build system. The exact
	// meaning is Type-dependent.
	Target string
	// Dir is the absolute path to the module's working directory (the directory
	// you would normally run the build command from).
	Dir string

	// TODO: add project and revision, because different modules might belong to
	// different projects
}

// New instantiates and sets up a Module for a given ModuleType
func New(moduleType Type, conf Config) (Module, error) {
	var manifestName string
	var moduleTarget string

	// Find root dir of module
	modulePath, err := filepath.Abs(conf.Path)
	if err != nil {
		return Module{}, err
	}

	// infer default module settings from type
	switch moduleType {
	case Bower:
		manifestName = "bower.json"
		break
	case Composer:
		manifestName = "composer.json"
		break
	case Golang:
		manifestName = ""
		moduleTarget = strings.TrimPrefix(modulePath, filepath.Join(os.Getenv("GOPATH"), "src")+"/")
		break
	case Maven:
		manifestName = "pom.xml"
		break
	case Nodejs:
		manifestName = "package.json"
		break
	case Ruby:
		manifestName = "Gemfile"
		break
	case SBT:
		manifestName = "build.sbt"
		break
	case VendoredArchives:
		manifestName = ""
		break
	}

	// trim manifest from path
	if filepath.Base(modulePath) == manifestName {
		modulePath = filepath.Dir(modulePath)
	}

	moduleName := conf.Name
	if moduleName == "" {
		moduleName = conf.Path
	}

	if moduleTarget == "" {
		moduleTarget = filepath.Join(modulePath, manifestName)
	}

	return Module{
		Name:   moduleName,
		Type:   moduleType,
		Target: moduleTarget,
		Dir:    modulePath,
	}, nil
}
