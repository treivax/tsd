// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// WithTimeout runs a test function with timeout
func WithTimeout(t *testing.T, duration time.Duration, fn func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-ctx.Done():
		t.Fatalf("Test timed out after %v", duration)
	}
}

// CreateTempTSDFile creates a temporary .tsd file for testing
func CreateTempTSDFile(t *testing.T, content string) string {
	t.Helper()

	// Create temp file with .tsd extension
	tmpFile, err := ioutil.TempFile("", "test-*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	path := tmpFile.Name()

	// Write content
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(path)
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(path)
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.Remove(path)
	})

	return path
}

// CreateTempTSDFileWithName creates a temp .tsd file with specific name
func CreateTempTSDFileWithName(t *testing.T, name, content string) string {
	t.Helper()

	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "tsd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create file with specific name
	path := filepath.Join(tmpDir, name)
	if !strings.HasSuffix(path, ".tsd") {
		path += ".tsd"
	}

	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return path
}

// CleanupTempFiles removes temporary test files
func CleanupTempFiles(t *testing.T, paths ...string) {
	t.Helper()

	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			t.Logf("Warning: Failed to cleanup %s: %v", path, err)
		}
	}
}

// SkipIfShort skips test if -short flag is set
func SkipIfShort(t *testing.T, reason string) {
	t.Helper()

	if testing.Short() {
		t.Skipf("Skipping test in short mode: %s", reason)
	}
}

// GetTestDataPath returns path to test data directory
func GetTestDataPath() string {
	projectRoot := getProjectRoot()
	return filepath.Join(projectRoot, "tests")
}

// GetFixturePath returns path to specific fixture
func GetFixturePath(category, name string) string {
	if !strings.HasSuffix(name, ".tsd") {
		name += ".tsd"
	}
	return filepath.Join(GetTestDataPath(), "fixtures", category, name)
}

// ReadTestFile reads a test file content
func ReadTestFile(t *testing.T, path string) string {
	t.Helper()

	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", path, err)
	}

	return string(content)
}

// WriteTestFile writes content to a test file
func WriteTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file %s: %v", path, err)
	}
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T, prefix string) string {
	t.Helper()

	tmpDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return tmpDir
}

// MustCreateDir creates a directory or fails the test
func MustCreateDir(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", path, err)
	}
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// CopyFile copies a file from src to dst
func CopyFile(t *testing.T, src, dst string) {
	t.Helper()

	content, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read source file %s: %v", src, err)
	}

	if err := ioutil.WriteFile(dst, content, 0644); err != nil {
		t.Fatalf("Failed to write destination file %s: %v", dst, err)
	}
}

// RequireFile fails the test if file doesn't exist
func RequireFile(t *testing.T, path string) {
	t.Helper()

	if !FileExists(path) {
		t.Fatalf("Required file does not exist: %s", path)
	}
}

// RequireDir fails the test if directory doesn't exist
func RequireDir(t *testing.T, path string) {
	t.Helper()

	if !DirExists(path) {
		t.Fatalf("Required directory does not exist: %s", path)
	}
}

// GetProjectRoot returns the project root directory
func GetProjectRoot(t *testing.T) string {
	t.Helper()
	return getProjectRoot()
}

// Retry retries a function until it succeeds or max attempts reached
func Retry(t *testing.T, attempts int, sleep time.Duration, fn func() error) error {
	t.Helper()

	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		if i < attempts-1 {
			time.Sleep(sleep)
		}
	}

	return err
}

// Eventually retries a condition until it's true or timeout
func Eventually(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Condition not met within %v: %s", timeout, message)
}

// CountFiles counts files matching pattern in directory
func CountFiles(t *testing.T, dir, pattern string) int {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		t.Fatalf("Failed to glob files: %v", err)
	}

	return len(matches)
}

// ListFiles lists all files in directory
func ListFiles(t *testing.T, dir string) []string {
	t.Helper()

	var files []string

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files
}

// SetEnv sets an environment variable for the test duration
func SetEnv(t *testing.T, key, value string) {
	t.Helper()

	oldValue, existed := os.LookupEnv(key)

	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set env var %s: %v", key, err)
	}

	// Restore on cleanup
	t.Cleanup(func() {
		if existed {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	})
}

// UnsetEnv unsets an environment variable for the test duration
func UnsetEnv(t *testing.T, key string) {
	t.Helper()

	oldValue, existed := os.LookupEnv(key)

	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed to unset env var %s: %v", key, err)
	}

	// Restore on cleanup
	t.Cleanup(func() {
		if existed {
			os.Setenv(key, oldValue)
		}
	})
}

// Chdir changes working directory for the test duration
func Chdir(t *testing.T, dir string) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory to %s: %v", dir, err)
	}

	// Restore on cleanup
	t.Cleanup(func() {
		os.Chdir(oldDir)
	})
}

// MeasureDuration measures execution duration of a function
func MeasureDuration(t *testing.T, name string, fn func()) time.Duration {
	t.Helper()

	start := time.Now()
	fn()
	duration := time.Since(start)

	t.Logf("%s took %v", name, duration)

	return duration
}

// AssertDuration fails if duration exceeds maximum
func AssertDuration(t *testing.T, duration, max time.Duration, operation string) {
	t.Helper()

	if duration > max {
		t.Errorf("%s took %v, expected at most %v", operation, duration, max)
	}
}
