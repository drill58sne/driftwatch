package digest_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/digest"
)

func sampleResult(host, name, status, output string) checker.CheckResult {
	return checker.CheckResult{
		Host:   host,
		Name:   name,
		Status: status,
		Output: output,
	}
}

func TestCompute_ReturnsDeterministicDigest(t *testing.T) {
	r := sampleResult("web-01", "uptime", "ok", "up 3 days")
	d1, err := digest.Compute(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d2, err := digest.Compute(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d1.Digest != d2.Digest {
		t.Errorf("expected identical digests, got %q and %q", d1.Digest, d2.Digest)
	}
}

func TestCompute_DifferentOutputProducesDifferentDigest(t *testing.T) {
	r1 := sampleResult("web-01", "uptime", "ok", "up 3 days")
	r2 := sampleResult("web-01", "uptime", "ok", "up 10 days")
	d1, _ := digest.Compute(r1)
	d2, _ := digest.Compute(r2)
	if d1.Digest == d2.Digest {
		t.Error("expected different digests for different outputs")
	}
}

func TestCompute_SetsHostAndName(t *testing.T) {
	r := sampleResult("db-01", "disk", "drift", "80%")
	d, err := digest.Compute(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Host != "db-01" {
		t.Errorf("expected host db-01, got %q", d.Host)
	}
	if d.Name != "disk" {
		t.Errorf("expected name disk, got %q", d.Name)
	}
}

func TestComputeAll_LengthMatchesInput(t *testing.T) {
	results := []checker.CheckResult{
		sampleResult("web-01", "uptime", "ok", "up"),
		sampleResult("web-02", "disk", "drift", "95%"),
	}
	digests, err := digest.ComputeAll(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(digests) != len(results) {
		t.Errorf("expected %d digests, got %d", len(results), len(digests))
	}
}

func TestChanged_ReturnsFalseWhenIdentical(t *testing.T) {
	r := sampleResult("web-01", "uptime", "ok", "up 3 days")
	prev, _ := digest.ComputeAll([]checker.CheckResult{r})
	curr, _ := digest.Compute(r)
	if digest.Changed(prev, curr) {
		t.Error("expected Changed=false for identical result")
	}
}

func TestChanged_ReturnsTrueWhenOutputDiffers(t *testing.T) {
	old := sampleResult("web-01", "uptime", "ok", "up 3 days")
	new_ := sampleResult("web-01", "uptime", "ok", "up 4 days")
	prev, _ := digest.ComputeAll([]checker.CheckResult{old})
	curr, _ := digest.Compute(new_)
	if !digest.Changed(prev, curr) {
		t.Error("expected Changed=true when output differs")
	}
}

func TestChanged_ReturnsTrueForUnknownHost(t *testing.T) {
	prev := []digest.Result{}
	curr, _ := digest.Compute(sampleResult("new-host", "check", "ok", "fine"))
	if !digest.Changed(prev, curr) {
		t.Error("expected Changed=true for unknown host")
	}
}
