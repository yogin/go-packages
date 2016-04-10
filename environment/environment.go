// Package environment allows the use of JSON configuration files for different environments
// It will look for an environment variable named GO_ENV to select an environment,
// it's also possible to pass the environment name via the Init() function
package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Environment struct {
	Name   string
	Config Config
}

type Config map[string]interface{}

var environments = map[string]string{
	"development": environmentFilePath("development"),
	"staging":     environmentFilePath("staging"),
	"production":  environmentFilePath("production"),
	"test":        environmentFilePath("test"),
}

var envDefault = "development"
var environment = Environment{}

// Register defines a new environment
func Register(name string) {
	if _, ok := environments[name]; !ok {
		environments[name] = environmentFilePath(name)
	}
}

// Init sets and loads the environment from a configuration file.
// It returns the Environment configuration
func Init(name ...string) Environment {
	var env string
	environment.Config = Config{}

	if len(name) > 0 {
		env = name[0]
		fmt.Printf("Warning: Overriding GO_ENV, setting senvironment to %s\n", env)
	} else {
		env = os.Getenv("GO_ENV")
		if env == "" {
			fmt.Printf("Warning: GO_ENV not set, setting environment to %s\n", envDefault)
			env = envDefault
		}
	}

	environment.Name = env

	if _, ok := environments[env]; !ok {
		panic(fmt.Sprintf("Error: Unsupported environment %v, Aborting!\n", env))
	}

	if err := load(environments[env]); err != nil {
		panic(fmt.Sprintf("Error: Failed reading configuration: %s\n", err))
	}

	return environment
}

// Get returns the environment
func Get() Environment {
	if environment.Name == "" {
		Init()
	}

	return environment
}

// Name returns the environment name
func Name() string {
	if environment.Name == "" {
		Init()
	}

	return environment.Name
}

// Get returns the value for a key
func (env Environment) Get(key string, defaultValue ...interface{}) interface{} {
	if env.Name == "" {
		panic(fmt.Sprintf("Error: Invalid environment: %+v\n", env))
	}

	var value interface{}

	if len(defaultValue) > 0 {
		value = defaultValue[0]
	}

	if _, ok := env.Config[key]; !ok {
		return value
	}

	return env.Config[key]
}

// Set adds/updates a key/value to the environment Config
func (env Environment) Set(key string, value interface{}) {
	if env.Name == "" {
		panic(fmt.Sprintf("Error: Invalid environment: %+v\n", env))
	}

	env.Config[key] = value
}

func load(path string) error {
	if environmentFileExists(path) == false {
		fmt.Printf("Warning: no environment configuration file found at %s\n", path)
		return nil
	}

	// Ignoring err here because I already check if the file exists
	// and ioutil doesn't check anything else: https://golang.org/src/io/ioutil/ioutil_test.go
	content, _ := ioutil.ReadFile(path)

	if err := json.Unmarshal(content, &environment.Config); err != nil {
		fmt.Printf("Error: failed parsing configuration file at %s : %s\n", path, err)
		return err
	}

	return nil
}

func environmentFileExists(path string) bool {
	exists := true

	if _, err := os.Stat(path); os.IsNotExist(err) {
		exists = false
	}

	return exists
}

func environmentFilePath(name string) string {
	return fmt.Sprintf("environments/%s.json", name)
}
