package inventory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Host represents a single remote server entry.
type Host struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	KeyPath  string `yaml:"key_path"`
	Tags     []string `yaml:"tags"`
}

// Inventory holds a list of hosts loaded from a file.
type Inventory struct {
	Hosts []Host `yaml:"hosts"`
}

// Load reads an inventory YAML file from the given path.
func Load(path string) (*Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading inventory file: %w", err)
	}

	var inv Inventory
	if err := yaml.Unmarshal(data, &inv); err != nil {
		return nil, fmt.Errorf("parsing inventory file: %w", err)
	}

	if err := inv.validate(); err != nil {
		return nil, err
	}

	return &inv, nil
}

// validate checks that all hosts have required fields.
func (inv *Inventory) validate() error {
	for i, h := range inv.Hosts {
		if h.Name == "" {
			return fmt.Errorf("host[%d]: name is required", i)
		}
		if h.Address == "" {
			return fmt.Errorf("host %q: address is required", h.Name)
		}
		if h.User == "" {
			return fmt.Errorf("host %q: user is required", h.Name)
		}
		if h.Port == 0 {
			inv.Hosts[i].Port = 22
		}
	}
	return nil
}

// FilterByTag returns hosts that have the given tag.
func (inv *Inventory) FilterByTag(tag string) []Host {
	var result []Host
	for _, h := range inv.Hosts {
		for _, t := range h.Tags {
			if t == tag {
				result = append(result, h)
				break
			}
		}
	}
	return result
}
