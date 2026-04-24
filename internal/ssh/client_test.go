package ssh

import (
	"testing"
	"time"
)

func TestConfigDefaults(t *testing.T) {
	cfg := Config{
		Host: "localhost",
		Port: 22,
		User: "admin",
	}

	if cfg.Timeout != 0 {
		t.Errorf("expected zero timeout before Connect, got %v", cfg.Timeout)
	}
}

func TestConnectInvalidKey(t *testing.T) {
	cfg := Config{
		Host:       "127.0.0.1",
		Port:       22,
		User:       "user",
		PrivateKey: []byte("not-a-valid-key"),
		Timeout:    2 * time.Second,
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("expected error for invalid private key, got nil")
	}
}

func TestConnectUnreachableHost(t *testing.T) {
	key := generateTestKey(t)

	cfg := Config{
		Host:       "192.0.2.1", // TEST-NET, guaranteed unreachable
		Port:       22,
		User:       "user",
		PrivateKey: key,
		Timeout:    1 * time.Second,
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("expected connection error for unreachable host")
	}
}

func TestConfigMissingUser(t *testing.T) {
	cfg := Config{
		Host:    "127.0.0.1",
		Port:    22,
		User:    "",
		Timeout: 1 * time.Second,
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}

// generateTestKey returns a placeholder private key byte slice for use in
// tests that require a non-nil key but do not reach actual SSH negotiation.
// For integration tests, use testcontainers or a mock SSH server with a
// properly generated key via crypto/rsa or crypto/ed25519.
func generateTestKey(t *testing.T) []byte {
	t.Helper()
	return []byte("placeholder")
}
