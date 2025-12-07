# Installation Guide

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [From Source](#from-source)
  - [Using Go Install](#using-go-install)
  - [Docker](#docker)
- [Verification](#verification)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required

- **Go 1.21 or higher**
  ```bash
  go version  # Should show 1.21 or higher
  ```

### Optional

- **Docker** (for containerized deployment)
- **Make** (for convenience commands)

## Installation Methods

### From Source

This is the recommended method for development and the most flexible installation option.

#### 1. Clone the Repository

```bash
git clone https://github.com/treivax/tsd.git
cd tsd
```

#### 2. Build the Binary

```bash
# Using Make (recommended)
make build

# Or using Go directly
go build -o bin/tsd ./cmd/tsd
```

#### 3. Verify Installation

```bash
./bin/tsd --version
```

#### 4. Install System-Wide (Optional)

```bash
# Linux/macOS
sudo cp bin/tsd /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)/bin
```

### Using Go Install

Install directly from the repository:

```bash
go install github.com/treivax/tsd/cmd/tsd@latest
```

The binary will be installed in `$GOPATH/bin` (typically `~/go/bin`).

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Docker

#### Pull Pre-built Image (if available)

```bash
docker pull treivax/tsd:latest
```

#### Build from Source

```bash
# From the project root
docker build -t tsd:local .
```

#### Run Container

```bash
# Run TSD server
docker run -p 8080:8080 tsd:local server

# Run TSD compiler (mount local directory)
docker run -v $(pwd):/workspace tsd:local /workspace/program.tsd
```

## Verification

### Basic Verification

```bash
# Check version
tsd --version

# Show help
tsd --help
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./rete
```

### Test with Example

```bash
# Create a simple test file
cat > test.tsd << 'EOF'
type Person(name: string, age: number)
action greet(name: string)

rule hello : {p: Person} / p.age >= 18 ==> greet(p.name)

Person(name: "Alice", age: 25)
EOF

# Run the program
tsd test.tsd
```

Expected output:
```
🎯 ACTION EXÉCUTÉE: greet("Alice")
```

## Configuration

### Binary Roles

TSD is a unified binary with multiple roles:

```bash
# Compiler/Runner (default)
tsd program.tsd

# Authentication management
tsd auth generate-key --output api-key.txt

# HTTP Server
tsd server --port 8080

# HTTP Client
tsd client --url http://localhost:8080 program.tsd
```

### Environment Variables

```bash
# Set log level
export TSD_LOG_LEVEL=debug

# Set server port
export TSD_PORT=8080

# Set authentication
export TSD_API_KEY=your-api-key-here

# Enable TLS
export TSD_TLS_CERT=/path/to/cert.pem
export TSD_TLS_KEY=/path/to/key.pem
```

### Configuration Files

Create a config file for persistent settings:

```yaml
# config.yaml
server:
  port: 8080
  host: 0.0.0.0
  
logging:
  level: info
  output: stdout
  
authentication:
  enabled: true
  key_file: /etc/tsd/api-key.txt
  
tls:
  enabled: false
  cert_file: /etc/tsd/cert.pem
  key_file: /etc/tsd/key.pem
```

Load configuration:

```bash
tsd server --config config.yaml
```

## Troubleshooting

### Build Issues

#### Go Version Too Old

```
Error: go version go1.20.x is too old
```

**Solution:** Upgrade Go to 1.21 or higher:

```bash
# Linux/macOS
go install golang.org/dl/go1.21.0@latest
go1.21.0 download
```

#### Missing Dependencies

```
Error: cannot find package ...
```

**Solution:** Download dependencies:

```bash
go mod download
go mod tidy
```

### Runtime Issues

#### Permission Denied

```
Error: permission denied when writing output
```

**Solution:** Check file permissions or run with appropriate privileges:

```bash
chmod +w output-directory/
# Or
sudo tsd program.tsd
```

#### Port Already in Use

```
Error: bind: address already in use
```

**Solution:** Use a different port or kill the process using the port:

```bash
# Find process using port 8080
lsof -i :8080
# Or
netstat -tulpn | grep 8080

# Kill the process
kill -9 <PID>

# Or use different port
tsd server --port 8081
```

#### TLS Certificate Issues

```
Error: tls: failed to verify certificate
```

**Solution:** Check certificate paths and validity:

```bash
# Verify certificate
openssl x509 -in cert.pem -text -noout

# Generate self-signed certificate for testing
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

### Testing Issues

#### Tests Fail Due to Race Conditions

```bash
# Run tests with race detector
go test -race ./...
```

#### Cleanup Test Artifacts

```bash
# Remove test cache
go clean -testcache

# Remove build artifacts
make clean
```

### Getting Help

1. **Check Documentation:**
   - [User Guide](USER_GUIDE.md)
   - [Tutorial](TUTORIAL.md)
   - [API Reference](API_REFERENCE.md)

2. **View Examples:**
   ```bash
   ls examples/
   ```

3. **Enable Debug Logging:**
   ```bash
   TSD_LOG_LEVEL=debug tsd program.tsd
   ```

4. **Report Issues:**
   - GitHub Issues: https://github.com/treivax/tsd/issues
   - Include: version, OS, error messages, minimal reproduction case

## Next Steps

After installation:

1. **Read the [Quick Start Guide](QUICK_START.md)** - Get up and running in 5 minutes
2. **Follow the [Tutorial](TUTORIAL.md)** - Learn TSD step by step
3. **Explore [Examples](../examples/)** - See real-world use cases
4. **Read the [User Guide](USER_GUIDE.md)** - Comprehensive feature documentation

## Uninstallation

### Remove Binary

```bash
# If installed system-wide
sudo rm /usr/local/bin/tsd

# If installed via go install
rm $(go env GOPATH)/bin/tsd
```

### Remove Source

```bash
# Remove cloned repository
rm -rf /path/to/tsd
```

### Remove Docker Images

```bash
docker rmi tsd:local
docker rmi treivax/tsd:latest
```

### Clean Go Cache

```bash
go clean -cache -modcache -i -r
```
