package model

import (
	"testing"
)

func TestMessage(t *testing.T) {
	expectedMessage := "test message"
	expectedServerCode := 200

	resp := Message(expectedMessage, expectedServerCode)

	if resp.Message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, resp.Message)
	}

	if resp.ServerCode != expectedServerCode {
		t.Errorf("Expected server code '%d', got '%d'", expectedServerCode, resp.ServerCode)
	}
}
