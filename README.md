# driftwatch

A lightweight CLI tool to detect and report config drift across remote servers via SSH.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Define your expected config state in a YAML file, then run `driftwatch` against your target servers:

```bash
driftwatch check --hosts hosts.txt --config baseline.yaml
```

**Example `baseline.yaml`:**

```yaml
files:
  - path: /etc/ssh/sshd_config
    contains: "PermitRootLogin no"
  - path: /etc/ntp.conf
    contains: "server time.example.com"
```

**Example output:**

```
[OK]    web-01  /etc/ssh/sshd_config
[DRIFT] web-02  /etc/ssh/sshd_config — expected "PermitRootLogin no" not found
[OK]    web-03  /etc/ntp.conf
[ERROR] web-04  /etc/ntp.conf — connection timed out
```

### Flags

| Flag | Description |
|------|-------------|
| `--hosts` | Path to a file containing hostnames or IPs (one per line) |
| `--config` | Path to the baseline config YAML |
| `--user` | SSH username (default: current user) |
| `--key` | Path to SSH private key (default: `~/.ssh/id_rsa`) |
| `--timeout` | SSH connection timeout (default: `10s`) |
| `--concurrency` | Number of hosts to check in parallel (default: `10`) |
| `--output` | Output format: `text` or `json` (default: `text`) |

---

## License

MIT © 2024 yourusername
