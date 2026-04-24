package ssh

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// startMockSSHServer spins up a local in-process SSH server for integration tests.
func startMockSSHServer(t *testing.T, handler func(ch ssh.Channel, reqs <-chan *ssh.Request)) (port int, hostKey ssh.Signer) {
	t.Helper()

	privKey, err := generateRSAKey()
	if err != nil {
		t.Fatalf("generating host key: %v", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(privKey)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listening: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		srvConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			return
		}
		defer srvConn.Close()
		go ssh.DiscardRequests(reqs)
		for newChan := range chans {
			if newChan.ChannelType() != "session" {
				_ = newChan.Reject(ssh.UnknownChannelType, "unsupported")
				continue
			}
			ch, chanReqs, _ := newChan.Accept()
			go handler(ch, chanReqs)
		}
	}()

	return ln.Addr().(*net.TCPAddr).Port, privKey
}

// echoHandler responds to exec requests with a fixed output.
func echoHandler(output string) func(ssh.Channel, <-chan *ssh.Request) {
	return func(ch ssh.Channel, reqs <-chan *ssh.Request) {
		defer ch.Close()
		for req := range reqs {
			if req.Type == "exec" {
				if req.WantReply {
					_ = req.Reply(true, nil)
				}
				_, _ = io.WriteString(ch, output)
				_, _ = ch.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
				return
			}
		}
	}
}

// generateRSAKey is a stub — real implementation uses crypto/rsa + x/crypto/ssh.
func generateRSAKey() (ssh.Signer, error) {
	return nil, fmt.Errorf("stub: use crypto/rsa.GenerateKey in real implementation")
}

func TestMockServer_Placeholder(t *testing.T) {
	// Placeholder to confirm the mock infrastructure compiles.
	_ = time.Second
	t.Log("mock SSH server scaffolding in place")
}
