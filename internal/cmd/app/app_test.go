package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

type Credentials struct {
	Username string
}

func TestSetUsername(t *testing.T) {
	home, _ := os.UserHomeDir()
	fullPath := filepath.Join(home, "passKeeper", ".config")
	cfgFilePath := filepath.Join(fullPath, cfgFile)

	// clean up after test
	defer func() {
		os.RemoveAll(fullPath)
	}()

	username := "testUser"
	err := SetUsername(username)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check if file exists
	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
		t.Errorf("File %s was not created", cfgFilePath)
	}

	// check the contents of the file
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		t.Errorf("Error reading file %s: %v", cfgFilePath, err)
	}

	creds := &Credentials{}
	err = yaml.Unmarshal(data, creds)
	if err != nil {
		t.Errorf("Error unmarshaling yaml: %v", err)
	}

	if creds.Username != username {
		t.Errorf("Expected username to be %s, got %s", username, creds.Username)
	}
}
func TestGetFileInfo(t *testing.T) {
	tests := []struct {
		meta      string
		want      []string
		expectErr bool
	}{
		{
			meta:      "fileName|extension|description",
			want:      []string{"", "fileName", "extension", "description"},
			expectErr: false,
		},
		{
			meta:      "invalidString",
			expectErr: true,
		},
		{
			meta:      "onlyTwo|parts",
			expectErr: true,
		},
		{
			meta:      "|empty|fields",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		got, err := GetFileInfo(tt.meta)
		if (err != nil) != tt.expectErr {
			t.Errorf("GetFileInfo() error = %v, expectErr %v", err, tt.expectErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetFileInfo() = %v, want %v", got, tt.want)
		}
	}
}
