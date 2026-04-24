package checker

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// CheckResult holds the outcome of a single config check on a remote host.
type CheckResult struct {
	Host    string
	Check   string
	Expected string
	Actual  string
	Drifted bool
	Err     error
}

// Check defines a named remote command whose output is compared to an expected value.
type Check struct {
	Name     string
	Command  string
	Expected string
}

// Runner executes checks against a remote host over an established SSH client.
type Runner struct {
	client *ssh.Client
	host   string
}

// NewRunner creates a Runner for the given SSH client and host label.
func NewRunner(client *ssh.Client, host string) *Runner {
	return &Runner{client: client, host: host}
}

// Run executes all provided checks and returns a slice of CheckResults.
func (r *Runner) Run(checks []Check) []CheckResult {
	results := make([]CheckResult, 0, len(checks))
	for _, c := range checks {
		results = append(results, r.runOne(c))
	}
	return results
}

func (r *Runner) runOne(c Check) CheckResult {
	result := CheckResult{
		Host:     r.host,
		Check:    c.Name,
		Expected: c.Expected,
	}

	session, err := r.client.NewSession()
	if err != nil {
		result.Err = fmt.Errorf("failed to open session: %w", err)
		return result
	}
	defer session.Close()

	out, err := session.Output(c.Command)
	if err != nil {
		result.Err = fmt.Errorf("command %q failed: %w", c.Command, err)
		return result
	}

	actual := strings.TrimSpace(string(out))
	result.Actual = actual
	result.Drifted = actual != strings.TrimSpace(c.Expected)
	return result
}
