package client

import (
	"encoding/base64"
	"os"
	secret "passKeeper/internal/models/secret"
	"reflect"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Test with a real file
	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if !FileExists(tempFile.Name()) {
		t.Errorf("expected file to exist")
	}

	// Test with a non-existent file
	if FileExists("nonexistentfile.xyz") {
		t.Errorf("expected file not to exist")
	}

	// Test with a real directory
	tempDir, err := os.MkdirTemp("", "exampledir")
	if err != nil {
		t.Fatalf("failed creating temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	if FileExists(tempDir) {
		t.Errorf("expected directory not to be recognized as file")
	}
}

func TestFileToBytes(t *testing.T) {
	// Set up a temporary file with known content
	expectedContent := []byte("Hello, world!")
	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(expectedContent)
	if err != nil {
		t.Fatalf("failed writing to temp file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		t.Fatalf("failed closing temp file: %v", err)
	}

	// Call FiletoBytes function
	resultContent, err := FiletoBytes(tempFile.Name())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Check the result
	if !reflect.DeepEqual(resultContent, expectedContent) {
		t.Errorf("expected %v, got %v", expectedContent, resultContent)
	}

	// Test with a non-existent file
	_, err = FiletoBytes("nonexistentfile.xyz")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestSecretRequestFromBytes(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		meta     string
		expected secret.SecretRequest
	}{
		{
			name: "case 1: typical byte slice",
			data: []byte("hello world"),
			meta: "meta1",
			expected: secret.SecretRequest{
				Type:     "ByteSlice",
				ByteData: base64.StdEncoding.EncodeToString([]byte("hello world")),
				Meta:     "meta1",
			},
		},
		{
			name: "case 2: empty byte slice",
			data: []byte(""),
			meta: "meta2",
			expected: secret.SecretRequest{
				Type:     "ByteSlice",
				ByteData: base64.StdEncoding.EncodeToString([]byte("")),
				Meta:     "meta2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SecretRequestFromBytes(tc.data, tc.meta)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
