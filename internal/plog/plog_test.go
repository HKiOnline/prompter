package plog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	logFile := "/tmp/test.log"
	plogger := New(logFile)

	if plogger == nil {
		t.Fatal("Expected non-nil Plogger")
	}

	if plogger.plogFile != logFile {
		t.Errorf("Expected plogFile '%s', got '%s'", logFile, plogger.plogFile)
	}
}

func TestWriteSingleMessage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write a single message
	plogger.Write(SERVER, "Test message")

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify the message was written
	if !strings.Contains(logContent, "Test message") {
		t.Errorf("Expected log to contain 'Test message', got: %s", logContent)
	}

	// Verify sender is included
	if !strings.Contains(logContent, "["+SERVER+"]") {
		t.Errorf("Expected log to contain '[%s]', got: %s", SERVER, logContent)
	}

	// Verify timestamp format (RFC3339)
	if !strings.Contains(logContent, "T") || !strings.Contains(logContent, ":") {
		t.Errorf("Expected log to contain RFC3339 timestamp, got: %s", logContent)
	}
}

func TestWriteMultipleMessages(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_multiple")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write multiple messages
	plogger.Write(CLIENT, "First message", "Second message", "Third message")

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify all messages are present
	if !strings.Contains(logContent, "First message") {
		t.Errorf("Expected log to contain 'First message', got: %s", logContent)
	}
	if !strings.Contains(logContent, "Second message") {
		t.Errorf("Expected log to contain 'Second message', got: %s", logContent)
	}
	if !strings.Contains(logContent, "Third message") {
		t.Errorf("Expected log to contain 'Third message', got: %s", logContent)
	}

	// Verify messages are separated by " :: "
	if !strings.Contains(logContent, " :: ") {
		t.Errorf("Expected messages to be separated by ' :: ', got: %s", logContent)
	}

	// Verify sender is CLIENT
	if !strings.Contains(logContent, "["+CLIENT+"]") {
		t.Errorf("Expected log to contain '[%s]', got: %s", CLIENT, logContent)
	}
}

func TestWriteAppendsToFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_append")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write first message
	plogger.Write(SERVER, "First log entry")

	// Write second message
	plogger.Write(SERVER, "Second log entry")

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify both messages are present
	if !strings.Contains(logContent, "First log entry") {
		t.Errorf("Expected log to contain 'First log entry', got: %s", logContent)
	}
	if !strings.Contains(logContent, "Second log entry") {
		t.Errorf("Expected log to contain 'Second log entry', got: %s", logContent)
	}

	// Verify there are two separate lines
	lines := strings.Split(strings.TrimSpace(logContent), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 log lines, got %d", len(lines))
	}
}

func TestWriteWithInvalidFilePath(t *testing.T) {
	// Use an invalid path (directory that doesn't exist and can't be created)
	logFile := "/nonexistent/directory/that/cannot/be/created/test.log"
	plogger := New(logFile)

	// This should not panic, but should handle the error gracefully
	// We can't easily verify stderr output, but we can verify it doesn't crash
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Write panicked with invalid file path: %v", r)
		}
	}()

	plogger.Write(SERVER, "This should not crash")
}

func TestWriteCreatesFileIfNotExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_create")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "new_log.log")

	// Verify file doesn't exist yet
	if _, err := os.Stat(logFile); err == nil {
		t.Fatal("Log file should not exist yet")
	}

	plogger := New(logFile)
	plogger.Write(SERVER, "Create new file")

	// Verify file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Expected log file to be created")
	}

	// Verify content was written
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Create new file") {
		t.Errorf("Expected log to contain 'Create new file', got: %s", string(content))
	}
}

func TestWriteFormatsTimestamp(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_timestamp")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	beforeWrite := time.Now()
	plogger.Write(SERVER, "Timestamp test")
	afterWrite := time.Now()

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Extract timestamp from log (it's at the beginning)
	parts := strings.Split(logContent, " ")
	if len(parts) < 2 {
		t.Fatalf("Expected log to have timestamp, got: %s", logContent)
	}

	timestampStr := parts[0] + " " + parts[1] // Date and time parts
	timestampStr = strings.TrimSpace(timestampStr)

	// Try to parse the timestamp
	parsedTime, err := time.Parse(time.RFC3339, strings.Split(logContent, " [")[0])
	if err != nil {
		t.Errorf("Failed to parse timestamp as RFC3339: %v, timestamp: %s", err, timestampStr)
	}

	// Verify timestamp is within reasonable range (within 1 minute of test execution)
	if parsedTime.Before(beforeWrite.Add(-1*time.Minute)) || parsedTime.After(afterWrite.Add(1*time.Minute)) {
		t.Errorf("Timestamp %v is not within expected range [%v, %v]", parsedTime, beforeWrite, afterWrite)
	}
}

func TestWriteWithDifferentSenders(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_senders")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write with SERVER sender
	plogger.Write(SERVER, "Server message")

	// Write with CLIENT sender
	plogger.Write(CLIENT, "Client message")

	// Write with custom sender
	plogger.Write("custom", "Custom message")

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify all senders are present
	if !strings.Contains(logContent, "["+SERVER+"]") {
		t.Errorf("Expected log to contain '[%s]'", SERVER)
	}
	if !strings.Contains(logContent, "["+CLIENT+"]") {
		t.Errorf("Expected log to contain '[%s]'", CLIENT)
	}
	if !strings.Contains(logContent, "[custom]") {
		t.Errorf("Expected log to contain '[custom]'")
	}
}

func TestWriteEmptyMessage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write with no messages
	plogger.Write(SERVER)

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify sender is still present
	if !strings.Contains(logContent, "["+SERVER+"]") {
		t.Errorf("Expected log to contain '[%s]', got: %s", SERVER, logContent)
	}

	// Verify line ends with newline
	if !strings.HasSuffix(logContent, "\n") {
		t.Errorf("Expected log line to end with newline")
	}
}

func TestWriteWithSpecialCharacters(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_special")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write with special characters
	specialMessage := "Message with special chars: !@#$%^&*(){}[]|\\:\";<>?,./~`"
	plogger.Write(SERVER, specialMessage)

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify special characters are preserved
	if !strings.Contains(logContent, specialMessage) {
		t.Errorf("Expected log to contain special characters message, got: %s", logContent)
	}
}

func TestWriteWithNewlinesInMessage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_newlines")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write with newlines in the message
	multilineMessage := "Line 1\nLine 2\nLine 3"
	plogger.Write(SERVER, multilineMessage)

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify newlines are preserved
	if !strings.Contains(logContent, multilineMessage) {
		t.Errorf("Expected log to contain multiline message, got: %s", logContent)
	}
}

func TestWriteConcurrency(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_concurrent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	// Write messages concurrently from multiple goroutines
	done := make(chan bool)
	numGoroutines := 10
	messagesPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < messagesPerGoroutine; j++ {
				plogger.Write(SERVER, "Goroutine", string(rune('0'+id)), "Message", string(rune('0'+j)))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Count the number of lines
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	expectedLines := numGoroutines * messagesPerGoroutine

	if len(lines) != expectedLines {
		t.Errorf("Expected %d log lines, got %d", expectedLines, len(lines))
	}

	// Verify each line has the expected format
	for _, line := range lines {
		if !strings.Contains(line, "["+SERVER+"]") {
			t.Errorf("Expected line to contain '[%s]', got: %s", SERVER, line)
		}
	}
}

func TestConstants(t *testing.T) {
	// Verify constants have expected values
	if CLIENT != "client" {
		t.Errorf("Expected CLIENT constant to be 'client', got '%s'", CLIENT)
	}

	if SERVER != "prompter" {
		t.Errorf("Expected SERVER constant to be 'prompter', got '%s'", SERVER)
	}
}

func TestWriteFilePermissions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_plog_perms")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	plogger := New(logFile)

	plogger.Write(SERVER, "Test permissions")

	// Check file permissions
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}

	// Verify file has expected permissions (0644)
	expectedPerms := os.FileMode(0644)
	actualPerms := fileInfo.Mode().Perm()

	if actualPerms != expectedPerms {
		t.Errorf("Expected file permissions %v, got %v", expectedPerms, actualPerms)
	}
}
