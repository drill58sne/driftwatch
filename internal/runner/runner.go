// Package runner orchestrates SSH connections, check execution, and reporting
// across all hosts defined in an inventory.
package runner

import (
	"fmt"
	"sync"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/inventory"
	"github.com/driftwatch/internal/reporter"
	"github.com/driftwatch/internal/ssh"
)

// Options controls runner behaviour.
type Options struct {
	Concurrency int
	Format      string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Concurrency: 5,
		Format:      "text",
	}
}

// Run connects to every host in inv, executes checks, and writes a report.
func Run(inv *inventory.Inventory, checks []checker.Check, opts Options, rep *reporter.Reporter) error {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 1
	}

	type result struct {
		res checker.CheckResult
		err error
	}

	sem := make(chan struct{}, opts.Concurrency)
	resultCh := make(chan result)
	var wg sync.WaitGroup

	for _, host := range inv.Hosts {
		wg.Add(1)
		go func(h inventory.Host) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			client, err := ssh.Connect(ssh.Config{
				Host:           h.Address,
				Port:           h.Port,
				User:           h.User,
				PrivateKeyPath: h.IdentityFile,
			})
			if err != nil {
				resultCh <- result{err: fmt.Errorf("host %s: connect: %w", h.Name, err)}
				return
			}
			defer client.Close()

			r := checker.NewRunner(client)
			for _, chk := range checks {
				res, err := r.Run(chk)
				resultCh <- result{res: res, err: err}
			}
		}(host)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []checker.CheckResult
	var errs []error
	for r := range resultCh {
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		results = append(results, r.res)
	}

	if err := rep.Write(results, opts.Format); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d host(s) failed: first error: %w", len(errs), errs[0])
	}
	return nil
}
