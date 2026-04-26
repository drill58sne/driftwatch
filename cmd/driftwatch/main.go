// Package main is the entry point for the driftwatch CLI tool.
// It wires together inventory loading, SSH connections, config checks,
// baseline comparison, alerting, and reporting into a single command.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/alert"
	"github.com/yourorg/driftwatch/internal/baseline"
	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/inventory"
	"github.com/yourorg/driftwatch/internal/reporter"
	"github.com/yourorg/driftwatch/internal/runner"
)

const version = "0.1.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("driftwatch", flag.ContinueOnError)

	var (
		inventoryFile = fs.String("inventory", "inventory.yaml", "Path to inventory YAML file")
		configFile    = fs.String("config", "driftwatch.yaml", "Path to driftwatch config file")
		baselineFile  = fs.String("baseline", "", "Path to baseline file (optional; enables drift comparison)")
		saveBaseline  = fs.Bool("save-baseline", false, "Save current results as the new baseline and exit")
		outputFmt     = fs.String("output", "text", "Output format: text or json")
		timeout       = fs.Duration("timeout", 30*time.Second, "SSH connection timeout per host")
		showVersion   = fs.Bool("version", false, "Print version and exit")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Printf("driftwatch %s\n", version)
		return nil
	}

	// Load inventory.
	inv, err := inventory.Load(*inventoryFile)
	if err != nil {
		return fmt.Errorf("loading inventory %q: %w", *inventoryFile, err)
	}

	// Load application config.
	cfg, err := config.Load(*configFile)
	if err != nil {
		return fmt.Errorf("loading config %q: %w", *configFile, err)
	}

	// Build runner options from flags and config.
	opts := runner.DefaultOptions()
	opts.Timeout = *timeout

	// Execute checks across all hosts.
	results, err := runner.Run(inv, cfg.Checks, opts)
	if err != nil {
		// Run returns a partial result set alongside any error; continue to
		// report what we have before propagating the error.
		fmt.Fprintf(os.Stderr, "warning: one or more hosts failed: %v\n", err)
	}

	// If requested, persist results as the new baseline and exit.
	if *saveBaseline {
		dest := *baselineFile
		if dest == "" {
			dest = "baseline.json"
		}
		if saveErr := baseline.Save(dest, results); saveErr != nil {
			return fmt.Errorf("saving baseline to %q: %w", dest, saveErr)
		}
		fmt.Printf("Baseline saved to %s (%d results)\n", dest, len(results))
		return nil
	}

	// Optionally compare against a baseline.
	if *baselineFile != "" {
		snap, loadErr := baseline.Load(*baselineFile)
		if loadErr != nil {
			return fmt.Errorf("loading baseline %q: %w", *baselineFile, loadErr)
		}
		results = baseline.Against(snap, results)
	}

	// Report results.
	rep := reporter.New()
	switch *outputFmt {
	case "json":
		if writeErr := rep.WriteJSON(results); writeErr != nil {
			return fmt.Errorf("writing JSON report: %w", writeErr)
		}
	default:
		if writeErr := rep.WriteText(results); writeErr != nil {
			return fmt.Errorf("writing text report: %w", writeErr)
		}
	}

	// Evaluate alert thresholds; exit non-zero when the level demands it.
	alertCfg := alert.DefaultConfig()
	alertCfg.WarnThreshold = cfg.Alert.WarnThreshold
	alertCfg.ErrorThreshold = cfg.Alert.ErrorThreshold

	al := alert.New(alertCfg)
	level := al.Evaluate(results)
	if level == alert.LevelError {
		os.Exit(2)
	}

	return err
}
