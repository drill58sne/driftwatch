package summarizer_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/summarizer"
)

func makeResults() []checker.Result {
	now := time.Now()
	return []checker.Result{
		{Host: "web-01", Name: "sshd", Drifted: false, CheckedAt: now},
		{Host: "web-01", Name: "motd", Drifted: true, CheckedAt: now},
		{Host: "web-02", Name: "sshd", Drifted: false, CheckedAt: now},
		{Host: "web-02", Name: "hosts", Error: errors.New("timeout"), CheckedAt: now},
	}
}

func TestCompute_TotalCounts(t *testing.T) {
	s := summarizer.Compute(makeResults())
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Drifted != 1 {
		t.Errorf("expected Drifted=1, got %d", s.Drifted)
	}
	if s.Clean != 2 {
		t.Errorf("expected Clean=2, got %d", s.Clean)
	}
	if s.Errored != 1 {
		t.Errorf("expected Errored=1, got %d", s.Errored)
	}
}

func TestCompute_HostCount(t *testing.T) {
	s := summarizer.Compute(makeResults())
	if len(s.Hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(s.Hosts))
	}
}

func TestCompute_DriftRate(t *testing.T) {
	s := summarizer.Compute(makeResults())
	for _, hs := range s.Hosts {
		if hs.Host == "web-01" {
			if hs.DriftRate != 0.5 {
				t.Errorf("expected DriftRate=0.5 for web-01, got %f", hs.DriftRate)
			}
		}
	}
}

func TestHasDrift_True(t *testing.T) {
	s := summarizer.Compute(makeResults())
	if !s.HasDrift() {
		t.Error("expected HasDrift()=true")
	}
}

func TestHasDrift_False(t *testing.T) {
	now := time.Now()
	results := []checker.Result{
		{Host: "db-01", Name: "sshd", Drifted: false, CheckedAt: now},
	}
	s := summarizer.Compute(results)
	if s.HasDrift() {
		t.Error("expected HasDrift()=false")
	}
}

func TestCompute_Empty(t *testing.T) {
	s := summarizer.Compute(nil)
	if s.Total != 0 || len(s.Hosts) != 0 {
		t.Error("expected empty summary for nil input")
	}
}
