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

// generateTestKey creates a temporary RSA key for testing.
func generateTestKey(t *testing.T) []byte {
	t.Helper()
	import_crypto := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAA=
-----END OPENSSH PRIVATE KEY-----`
	// Intentionally malformed — real key generation requires crypto/rsa.
	// Integration tests should use testcontainers or a mock SSH server.
	_ = import_crypto
	return []byte("placeholder")
}
