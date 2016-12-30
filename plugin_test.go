package main

import (
	"github.com/stretchr/testify/assert"

	"os"
	"os/user"
	"testing"
)

func TestMissingAllConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingSSHConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"example.com"},
			Username: "ubuntu",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingSourceConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"example.com"},
			Username: "ubuntu",
			Port:     "443",
			Password: "1234",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestTrimElement(t *testing.T) {
	var input, result []string

	input = []string{"1", "     ", "3"}
	result = []string{"1", "3"}

	assert.Equal(t, result, trimPath(input))

	input = []string{"1", "2"}
	result = []string{"1", "2"}

	assert.Equal(t, result, trimPath(input))
}

func TestSCPFile(t *testing.T) {
	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     "22",
			KeyPath:  "tests/.ssh/id_rsa",
			Source:   []string{"tests/a.txt", "tests/b.txt"},
			Target:   []string{u.HomeDir + "/test"},
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(u.HomeDir + "/test/tests/a.txt"); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(u.HomeDir + "/test/tests/b.txt"); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	// Test -rm flag
	plugin.Config.Source = []string{"tests/a.txt"}
	plugin.Config.Remove = true

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(u.HomeDir + "/test/tests/b.txt"); os.IsExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestIncorrectPassword(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     "22",
			Password: "1234",
			Source:   []string{"tests/a.txt", "tests/b.txt"},
			Target:   []string{"/home"},
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestNoPermissionCreateFolder(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     "22",
			KeyPath:  "tests/.ssh/id_rsa",
			Source:   []string{"tests/a.txt", "tests/b.txt"},
			Target:   []string{"/etc/test"},
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSourceNotFound(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     "22",
			KeyPath:  "tests/.ssh/id_rsa",
			Source:   []string{"tests/aa.txt", "tests/b.txt"},
			Target:   []string{"/test"},
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}
