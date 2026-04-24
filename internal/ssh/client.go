package ssh

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// Config holds SSH connection parameters.
type Config struct {
	Host       string
	Port       int
	User       string
	PrivateKey []byte
	Timeout    time.Duration
}

// Client wraps an SSH client connection.
type Client struct {
	conn *ssh.Client
}

// Connect establishes an SSH connection using the provided config.
func Connect(cfg Config) (*Client, error) {
	signer, err := ssh.ParsePrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	sshCfg := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: replace with known_hosts
		Timeout:         cfg.Timeout,
	}

	addr := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return nil, fmt.Errorf("dialing %s: %w", addr, err)
	}

	return &Client{conn: conn}, nil
}

// RunCommand executes a command on the remote host and returns combined output.
func (c *Client) RunCommand(cmd string) (string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("creating session: %w", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(out), fmt.Errorf("running command %q: %w", cmd, err)
	}

	return string(out), nil
}

// Close terminates the SSH connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
