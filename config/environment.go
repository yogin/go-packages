package config

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

// Init sets and loads the environment from a configuration file
func Init(name ...string) {
	var env string

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
}

// GetEnvironment returns the environment
func Get() Environment {
	if environment.Name == "" {
		Init()
	}

	return environment
}

// GetName returns the environment name
func Name() string {
	if environment.Name == "" {
		Init()
	}

	return environment.Name
}

// Get returns the value for a key
func (c Config) Get(key string, defaultValue ...interface{}) interface{} {
	var value interface{}

	if len(defaultValue) > 0 {
		value = defaultValue[0]
	}

	if _, ok := c[key]; !ok {
		return value
	}

	return c[key]
}

func load(path string) error {
	if environmentFileExists(path) == false {
		fmt.Printf("Warning: no environment configuration file found at %s\n", path)
		return nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: failed reading configuration file at %s : %s\n", path, err)
		return err
	}

	if err := json.Unmarshal(content, &environment.Config); err != nil {
		fmt.Printf("Error: failed parsing configuration file at %s : %s\n", path, err)
		return err
	}

	return nil
}

func environmentFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func environmentFilePath(name string) string {
	return fmt.Sprintf("config/environments/%s.json", name)
}
