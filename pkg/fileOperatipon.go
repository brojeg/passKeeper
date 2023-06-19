package client

import (
	"encoding/base64"
	"fmt"
	"os"
	secret "passKeeper/internal/models/secret"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FiletoBytes(path string) ([]byte, error) {

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return fileData, nil

}

func SecretRequestFromBytes(data []byte, meta string) secret.SecretRequest {

	base64Data := base64.StdEncoding.EncodeToString(data)

	converteddata := secret.ByteSlice(base64Data)

	request := secret.SecretRequest{
		Type:     "ByteSlice",
		ByteData: string(converteddata),
		Meta:     meta,
	}

	return request
}
