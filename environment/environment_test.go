package environment

import (
	"os"
	"testing"
)

func resetEnvironment() {
	environment = Environment{}
}

func TestInit_WithoutEnv(t *testing.T) {
	resetEnvironment()
	Init()

	if environment.Name != envDefault {
		t.Errorf("got: %s, expected: %s", environment.Name, envDefault)
	}
}

func TestInit_WithValidEnv(t *testing.T) {
	resetEnvironment()

	old := os.Getenv("GO_ENV")
	defer func() { os.Setenv("GO_ENV", old) }()

	os.Setenv("GO_ENV", "test")
	Init()

	if environment.Name != "test" {
		t.Errorf("got: %s, expected: %s", environment.Name, "test")
	}
}

func TestInit_WithInvalidEnv(t *testing.T) {
	resetEnvironment()

	old := os.Getenv("GO_ENV")
	defer func() { os.Setenv("GO_ENV", old) }()

	os.Setenv("GO_ENV", "panic_env")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("got: %s, expected: panic!", environment.Name)
		}
	}()

	Init()
}

func ExampleF_Init() {
	resetEnvironment()
	Init()
	// Output:
	// Warning: GO_ENV not set, setting environment to development
	// Warning: no environment configuration file found at environments/development.json
}

func ExampleF_InitValid() {
	resetEnvironment()
	Init("test")
	// Output:
	// Warning: Overriding GO_ENV, setting senvironment to test
	// Warning: no environment configuration file found at environments/test.json
}

func TestGetDefault(t *testing.T) {
	resetEnvironment()
	Init()
	env := Get()

	if env.Name != envDefault {
		t.Errorf("got: %s, expected: %s", env.Name, envDefault)
	}
}

func TestGetWithValue(t *testing.T) {
	resetEnvironment()
	Init("test")
	env := Get()

	if env.Name != "test" {
		t.Errorf("got: %s, expected: %s", env.Name, "test")
	}
}

func TestGetWithoutInit(t *testing.T) {
	resetEnvironment()
	env := Get()

	if env.Name != "development" {
		t.Errorf("got: %s, expected %s", env.Name, "development")
	}
}

func TestNameDefault(t *testing.T) {
	resetEnvironment()
	Init()
	name := Name()

	if name != envDefault {
		t.Errorf("got: %s, expected: %s", name, envDefault)
	}
}

func TestNameWithValue(t *testing.T) {
	resetEnvironment()
	Init("test")
	name := Name()

	if name != "test" {
		t.Errorf("got: %s, expected: %s", name, "test")
	}
}

func TestNameWithoutInit(t *testing.T) {
	resetEnvironment()
	name := Name()

	if name != envDefault {
		t.Errorf("got: %s, expected: %s", name, envDefault)
	}
}

func TestRegister(t *testing.T) {
	resetEnvironment()
	Register("new_environment")

	if _, ok := environments["new_environment"]; !ok {
		t.Errorf("environment 'new_environment' not found")
	}
}

func TestConfigGet(t *testing.T) {
	resetEnvironment()
	Init()
	env := Get()
	env.Config["world"] = "hello"
	value := env.Get("world")

	if value != "hello" {
		t.Errorf("got: %v, expected: hello", value)
	}
}

func TestConfigGetDefault(t *testing.T) {
	resetEnvironment()
	Init()
	env := Get()
	value := env.Get("world", 42)

	if value != 42 {
		t.Errorf("got: %v, expected: 42", value)
	}
}

func TestConfigGetWithInvalidEnv(t *testing.T) {
	env := Environment{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("got: %+v, expected: panic!", env)
		}
	}()

	env.Get("hello")
}

func TestConfigSet(t *testing.T) {
	resetEnvironment()
	Init()
	env := Get()
	env.Set("hello", "world")

	if _, ok := env.Config["hello"]; !ok {
		t.Errorf("expecting to set config 'hello'='world'")
	}
}

func TestConfigSetWithInvalidEnv(t *testing.T) {
	env := Environment{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("got: %+v, expected: panic!", env)
		}
	}()

	env.Set("hello", "world")
}

func TestEnvironmentFileExists(t *testing.T) {
	path := environmentFilePath("sample_valid")
	v := environmentFileExists(path)

	if v != true {
		t.Errorf("got: %v, expected: %v", v, true)
	}
}

func TestEnvironmentFileNotExists(t *testing.T) {
	path := environmentFilePath("not_exist")
	v := environmentFileExists(path)

	if v != false {
		t.Errorf("got: %v, expected: %v", v, true)
	}
}

func TestLoadValid(t *testing.T) {
	resetEnvironment()
	path := environmentFilePath("sample_valid")

	if err := load(path); err != nil {
		t.Errorf("expected to load, got: %s", err)
	}
}

func TestLoadInValid(t *testing.T) {
	resetEnvironment()
	path := environmentFilePath("sample_invalid")

	if err := load(path); err == nil {
		t.Errorf("expected to fail loading, got: %+v", environment.Config)
	}
}

func TestInitWithInvalidConfig(t *testing.T) {
	Register("sample_invalid")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("got: %+v, expected: panic!", environment.Config)
		}
	}()

	Init("sample_invalid")
}
