# Atlas

Atlas is a small, opinionated deployment tool for Linux servers using **rsync + systemd**.

It deploys files into versioned release directories, atomically updates a `current` symlink, installs or updates a systemd service, and restarts the application.

## Features
- File-based deploys via rsync
- Versioned releases with rollback support
- systemd unit generation and management
- Language-agnostic (Go, Python, etc.)

## Example config

```yaml
app:
  name: golang-server

deploy:
  local-dir: ./app
  user: deploy
  host: example.com

run:
  user: app
  command: ["./app"]