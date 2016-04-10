package environment

import (
	"os"
	"testing"
)

func TestInit_WithoutEnv(t *testing.T) {
	Init()

	if environment.Name != envDefault {
		t.Errorf("got: %s, expected: %s", environment.Name, envDefault)
	}
}

func TestInit_WithValidEnv(t *testing.T) {
	old := os.Getenv("GO_ENV")
	defer func() { os.Setenv("GO_ENV", old) }()

	os.Setenv("GO_ENV", "test")
	Init()

	if environment.Name != "test" {
		t.Errorf("got: %s, expected: %s", environment.Name, "test")
	}
}

func TestInit_WithInvalidEnv(t *testing.T) {
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
	Init()
	// Output:
	// Warning: GO_ENV not set, setting environment to development
	// Warning: no environment configuration file found at environments/development.json
}

func ExampleF_InitValid() {
	Init("test")
	// Output:
	// Warning: Overriding GO_ENV, setting senvironment to test
	// Warning: no environment configuration file found at environments/test.json
}

func TestGetDefault(t *testing.T) {
	Init()
	env := Get()

	if env.Name != envDefault {
		t.Errorf("got: %s, expected: %s", env.Name, envDefault)
	}
}
func TestGetWithValue(t *testing.T) {
	Init("test")
	env := Get()

	if env.Name != "test" {
		t.Errorf("got: %s, expected: %s", env.Name, "test")
	}
}

func TestNameDefault(t *testing.T) {
	Init()
	name := Name()

	if name != envDefault {
		t.Errorf("got: %s, expected: %s", name, envDefault)
	}
}

func TestNameWithValue(t *testing.T) {
	Init("test")
	name := Name()

	if name != "test" {
		t.Errorf("got: %s, expected: %s", name, "test")
	}
}

func TestRegister(t *testing.T) {
	t.Skip("TODO")
}

func TestConfigGet(t *testing.T) {
	t.Skip("TODO")
}
