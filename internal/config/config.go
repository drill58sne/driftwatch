// Package config handles loading and validating driftwatch configuration,
// including server inventory and the checks to run against each host.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Check defines a single drift check to perform on a remote host.
type Check struct {
	// Name is a human-readable identifier for the check.
	Name string `yaml:"name"`
	// Command is the shell command to run on the remote host.
	Command string `yaml:"command"`
	// Expected is the expected output of the command (exact string match).
	Expected string `yaml:"expected"`
}

// Host represents a remote server to inspect.
type Host struct {
	// Address is the hostname or IP of the remote server.
	Address string `yaml:"address"`
	// Port is the SSH port; defaults to 22 if not specified.
	Port int `yaml:"port"`
	// User is the SSH username to authenticate as.
	User string `yaml:"user"`
	// IdentityFile is the path to the private key used for authentication.
	// If empty, the default SSH agent or ~/.ssh/id_rsa is used.
	IdentityFile string `yaml:"identity_file"`
}

// Config is the top-level structure representing a driftwatch configuration file.
type Config struct {
	// Hosts is the list of remote servers to check.
	Hosts []Host `yaml:"hosts"`
	// Checks is the list of drift checks to run against every host.
	Checks []Check `yaml:"checks"`
}

// Load reads and parses a YAML configuration file at the given path.
// It returns an error if the file cannot be read or the YAML is invalid.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks that the configuration contains at least one host and one
// check, and that each entry has the required fields populated.
func (c *Config) validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("at least one host must be defined")
	}
	if len(c.Checks) == 0 {
		return fmt.Errorf("at least one check must be defined")
	}

	for i, h := range c.Hosts {
		if h.Address == "" {
			return fmt.Errorf("host[%d]: address is required", i)
		}
		if h.User == "" {
			return fmt.Errorf("host[%d] (%s): user is required", i, h.Address)
		}
		if h.Port == 0 {
			c.Hosts[i].Port = 22
		}
	}

	for i, ch := range c.Checks {
		if ch.Name == "" {
			return fmt.Errorf("check[%d]: name is required", i)
		}
		if ch.Command == "" {
			return fmt.Errorf("check[%d] (%s): command is required", i, ch.Name)
		}
	}

	return nil
}
