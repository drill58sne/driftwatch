package plugin_test

import (
	"errors"
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/plugin"
)

func makePlugin(name string) plugin.Plugin {
	return plugin.Plugin{
		Name:    name,
		Version: "1.0.0",
		Check: func(host string) (checker.CheckResult, error) {
			return checker.CheckResult{Name: name, Output: "ok"}, nil
		},
	}
}

func TestRegister_Success(t *testing.T) {
	reg := plugin.New()
	if err := reg.Register(makePlugin("p1")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(makePlugin("dup"))
	if err := reg.Register(makePlugin("dup")); err == nil {
		t.Fatal("expected error for duplicate plugin name")
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	reg := plugin.New()
	p := plugin.Plugin{Name: "", Version: "1", Check: func(h string) (checker.CheckResult, error) { return checker.CheckResult{}, nil }}
	if err := reg.Register(p); err == nil {
		t.Fatal("expected error for empty plugin name")
	}
}

func TestRegister_NilCheck_ReturnsError(t *testing.T) {
	reg := plugin.New()
	p := plugin.Plugin{Name: "nil-check", Version: "1", Check: nil}
	if err := reg.Register(p); err == nil {
		t.Fatal("expected error for nil Check function")
	}
}

func TestGet_NotFound_ReturnsError(t *testing.T) {
	reg := plugin.New()
	if _, err := reg.Get("missing"); err == nil {
		t.Fatal("expected error for missing plugin")
	}
}

func TestList_ReturnsRegisteredNames(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(makePlugin("a"))
	_ = reg.Register(makePlugin("b"))
	names := reg.List()
	if len(names) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(names))
	}
}

func TestUnregister_RemovesPlugin(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(makePlugin("remove-me"))
	reg.Unregister("remove-me")
	if _, err := reg.Get("remove-me"); err == nil {
		t.Fatal("expected plugin to be removed")
	}
}

func TestRegister_CheckFnError_IsPreserved(t *testing.T) {
	reg := plugin.New()
	want := errors.New("check failed")
	_ = reg.Register(plugin.Plugin{
		Name:    "failing",
		Version: "1",
		Check:   func(h string) (checker.CheckResult, error) { return checker.CheckResult{}, want },
	})
	p, _ := reg.Get("failing")
	_, err := p.Check("host1")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
