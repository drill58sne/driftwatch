package inventory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/inventory"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "inventory.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoad_ValidInventory(t *testing.T) {
	path := writeTemp(t, `
hosts:
  - name: web-01
    address: 192.168.1.10
    user: admin
    port: 22
    tags: [web, prod]
`)
	inv, err := inventory.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inv.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(inv.Hosts))
	}
	if inv.Hosts[0].Name != "web-01" {
		t.Errorf("expected name web-01, got %s", inv.Hosts[0].Name)
	}
}

func TestLoad_DefaultPort(t *testing.T) {
	path := writeTemp(t, `
hosts:
  - name: db-01
    address: 10.0.0.5
    user: ubuntu
`)
	inv, err := inventory.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.Hosts[0].Port != 22 {
		t.Errorf("expected default port 22, got %d", inv.Hosts[0].Port)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTemp(t, `
hosts:
  - name: bad-host
    user: root
`)
	_, err := inventory.Load(path)
	if err == nil {
		t.Fatal("expected error for missing address, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := inventory.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFilterByTag(t *testing.T) {
	path := writeTemp(t, `
hosts:
  - name: web-01
    address: 1.1.1.1
    user: admin
    tags: [web, prod]
  - name: db-01
    address: 1.1.1.2
    user: admin
    tags: [db, prod]
`)
	inv, err := inventory.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	webHosts := inv.FilterByTag("web")
	if len(webHosts) != 1 || webHosts[0].Name != "web-01" {
		t.Errorf("expected 1 web host, got %v", webHosts)
	}
	prodHosts := inv.FilterByTag("prod")
	if len(prodHosts) != 2 {
		t.Errorf("expected 2 prod hosts, got %d", len(prodHosts))
	}
}
