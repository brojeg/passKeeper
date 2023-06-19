package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type testByteConvertible struct {
	data ByteSlice
}

func (t *testByteConvertible) ToBytes() (ByteSlice, error) {
	return t.data, nil
}

func (t *testByteConvertible) FromBytes(b ByteSlice) error {
	t.data = b
	return nil
}
func TestNewSecret(t *testing.T) {
	userID := uint(1)
	secretType := "TestType"
	value := &testByteConvertible{data: []byte("test data")}
	meta := "TestMeta"

	secret, err := NewSecret(userID, secretType, value, meta)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if secret.UserID != userID {
		t.Errorf("Expected userID '%d', got '%d'", userID, secret.UserID)
	}

	if secret.SecretType != secretType {
		t.Errorf("Expected secretType '%s', got '%s'", secretType, secret.SecretType)
	}

	if string(secret.Value) != string(value.data) {
		t.Errorf("Expected value '%s', got '%s'", string(value.data), string(secret.Value))
	}

	if secret.Metadata != meta {
		t.Errorf("Expected meta '%s', got '%s'", meta, secret.Metadata)
	}
}
func TestGetSecretFromRequest(t *testing.T) {
	testCases := []struct {
		name          string
		req           SecretRequest
		user          uint
		expectedValue ByteConvertible
		expectedErr   error
	}{
		{
			name: "valid KeyValue request",
			req: SecretRequest{
				Type: "KeyValue",
				Data: json.RawMessage(`{"key":"mykey","value":"myvalue"}`),
			},
			user:          uint(1),
			expectedValue: &KeyValue{Key: "mykey", Value: "myvalue"},
			expectedErr:   nil,
		},
		{
			name: "invalid type",
			req: SecretRequest{
				Type: "InvalidType",
			},
			user:          uint(1),
			expectedValue: nil,
			expectedErr:   fmt.Errorf("invalid type: InvalidType"),
		},
		// Add test cases for "Text", "CreditCard", and "ByteSlice" as well.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualValue, actualErr := GetSecretFromRequest(tc.req, tc.user)
			if !reflect.DeepEqual(tc.expectedValue, actualValue) {
				t.Errorf("Expected value %+v, but got %+v", tc.expectedValue, actualValue)
			}
			if !reflect.DeepEqual(tc.expectedErr, actualErr) {
				t.Errorf("Expected error %v, but got %v", tc.expectedErr, actualErr)
			}
		})
	}
}

func TestGetDecodedSecrets(t *testing.T) {
	testCases := []struct {
		name            string
		secrets         []Secret
		expectedDecoded []DecodedSecret
		expectedErr     error
	}{
		{
			name: "valid KeyValue secret",
			secrets: []Secret{
				{
					ID:         uint(1),
					UserID:     uint(1),
					Value:      ByteSlice(`{"key":"mykey","value":"myvalue"}`),
					SecretType: "KeyValue",
					Metadata:   "test",
				},
			},
			expectedDecoded: []DecodedSecret{
				{
					ID:       uint(1),
					UserID:   uint(1),
					Value:    &KeyValue{Key: "mykey", Value: "myvalue"},
					Metadata: "test",
				},
			},
			expectedErr: nil,
		},
		{
			name: "valid Text secret",
			secrets: []Secret{
				{
					ID:         uint(1),
					UserID:     uint(1),
					Value:      ByteSlice(`{"value":"myvalue"}`),
					SecretType: "Text",
					Metadata:   "test",
				},
			},
			expectedDecoded: []DecodedSecret{
				{
					ID:       uint(1),
					UserID:   uint(1),
					Value:    &Text{Value: "myvalue"},
					Metadata: "test",
				},
			},
			expectedErr: nil,
		},
		{
			name: "valid CreditCard secret",
			secrets: []Secret{
				{
					ID:         uint(1),
					UserID:     uint(1),
					Value:      ByteSlice(`{"number":"1234567890123456","Expiration":"05/23","cvv":"123"}`),
					SecretType: "CreditCard",
					Metadata:   "test",
				},
			},
			expectedDecoded: []DecodedSecret{
				{
					ID:       uint(1),
					UserID:   uint(1),
					Value:    &CreditCard{Number: "1234567890123456", Expiration: "05/23", CVV: "123"},
					Metadata: "test",
				},
			},
			expectedErr: nil,
		},
		{
			name: "unknown secret type",
			secrets: []Secret{
				{
					ID:         uint(1),
					UserID:     uint(1),
					Value:      ByteSlice(`{}`),
					SecretType: "UnknownType",
					Metadata:   "test",
				},
			},
			expectedDecoded: nil,
			expectedErr:     fmt.Errorf("unknown secret type: UnknownType"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualDecoded, actualErr := GetDecodedSecrets(tc.secrets)
			if !reflect.DeepEqual(tc.expectedDecoded, actualDecoded) {
				t.Errorf("Expected decoded secrets %+v, but got %+v", tc.expectedDecoded, actualDecoded)
			}
			if !reflect.DeepEqual(tc.expectedErr, actualErr) {
				t.Errorf("Expected error %v, but got %v", tc.expectedErr, actualErr)
			}
		})
	}
}

func TestValueToString(t *testing.T) {
	testCases := []struct {
		name           string
		secret         DecodedSecret
		expectedOutput string
	}{
		{
			name: "KeyValue",
			secret: DecodedSecret{
				Value: &KeyValue{Key: "mykey", Value: "myvalue"},
			},
			expectedOutput: "Key: mykey,\nValue: myvalue",
		},
		{
			name: "Text",
			secret: DecodedSecret{
				Value: &Text{Value: "myvalue"},
			},
			expectedOutput: "myvalue",
		},
		{
			name: "CreditCard",
			secret: DecodedSecret{
				Value: &CreditCard{Number: "1234567890123456", Expiration: "05/23", CVV: "123", Cardholder: "John Doe"},
			},
			expectedOutput: "Number: 1234567890123456,\n Expiration: 05/23,\n CVV: 123,\n Cardholder: John Doe",
		},
		{
			name: "ByteSlice with valid metadata",
			secret: func() DecodedSecret {
				b := ByteSlice([]byte("Hello"))
				return DecodedSecret{
					Value:    &b,
					Metadata: "myFile|txt|This is a text file",
				}
			}(),
			expectedOutput: "This is a text file",
		},
		{
			name: "ByteSlice with invalid metadata",
			secret: func() DecodedSecret {
				b := ByteSlice([]byte("Hello"))
				return DecodedSecret{
					Value:    &b,
					Metadata: "invalid metadata",
				}
			}(),
			expectedOutput: "*Binary data*",
		},
		{
			name: "Unknown Value Type",
			secret: DecodedSecret{
				Value: struct{}{}, // An empty struct represents an unknown value type.
			},
			expectedOutput: "Unknown Value Type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.secret.ValueToString()
			if actualOutput != tc.expectedOutput {
				t.Errorf("Expected output %s, but got %s", tc.expectedOutput, actualOutput)
			}
		})
	}
}
