package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	secret "passKeeper/internal/models/secret"
	"strings"
	"testing"
)

func TestSendJSONRequest(t *testing.T) {
	var payload = struct {
		Name string `json:"name"`
	}{
		Name: "TestName",
	}

	token := "Bearer token"

	tests := []struct {
		name     string
		method   string
		host     string
		endpoint string
		token    string
		payload  interface{}
		status   int
		body     string
		hasErr   bool
	}{
		{
			name:     "successful request",
			method:   http.MethodPost,
			host:     "", // the host will be replaced with the test server host
			endpoint: "/test",
			token:    token,
			payload:  &payload,
			status:   http.StatusOK,
			body:     `{"message":"success"}`,
			hasErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			resp, err := sendJSONRequest(client, tt.method, tt.host, tt.endpoint, tt.token, tt.payload)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}

			if !tt.hasErr && string(resp) != tt.body+"\n" {
				t.Fatalf("expected body %v, got %v", tt.body, string(resp))
			}
		})
	}
}

func TestSendRegisterRequest(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		login    string
		password string
		status   int
		body     string
		hasErr   bool
	}{
		{
			name:     "successful register",
			host:     "",
			login:    "testuser",
			password: "testpass",
			status:   http.StatusOK,
			body:     `{"login":"testuser","password":"testpass"}`,
			hasErr:   false,
		},
		{
			name:     "failed register due to missing data",
			host:     "",
			login:    "",
			password: "",
			status:   http.StatusBadRequest,
			body:     "",
			hasErr:   true,
		},
		{
			name:     "client error",
			host:     "",
			login:    "testuser",
			password: "testpass",
			status:   http.StatusBadRequest,
			body:     `{"message":"bad request"}`,
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			_, err := SendRegisterRequest(client, tt.host, tt.login, tt.password)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}

func TestSendGetSecretList(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		token  string
		status int
		body   string
		hasErr bool
	}{
		{
			name:   "successful get secret list",
			host:   "", // the host will be replaced with the test server host
			token:  "testToken",
			status: http.StatusOK,
			body: `[
				{
					"ID": 1,
					"UserID": 1,
					"Value": "dGVzdFNlY3JldA==",
					"SecretType": "testType",
					"Metadata": "testMetadata"
				}
			]`,
			hasErr: false,
		},
		{
			name:   "failed get secret list due to missing token",
			host:   "",
			token:  "",
			status: http.StatusUnauthorized,
			body:   "",
			hasErr: true,
		},
		{
			name:   "client error",
			host:   "",
			token:  "testToken",
			status: http.StatusBadRequest,
			body:   `{"message":"bad request"}`,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			_, err := SendGetSecretList(client, tt.host, tt.token)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}

func TestSendLoginRequest(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		login    string
		password string
		status   int
		body     string
		hasErr   bool
	}{
		{
			name:     "successful login",
			host:     "", // the host will be replaced with the test server host
			login:    "testuser",
			password: "testpass",
			status:   http.StatusOK,
			body:     `"testToken"`,
			hasErr:   false,
		},
		{
			name:     "failed login due to missing data",
			host:     "",
			login:    "",
			password: "",
			status:   http.StatusBadRequest,
			body:     "",
			hasErr:   true,
		},
		{
			name:     "client error",
			host:     "",
			login:    "testuser",
			password: "testpass",
			status:   http.StatusBadRequest,
			body:     `{"message":"bad request"}`,
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			_, err := SendLoginRequest(client, tt.host, tt.login, tt.password)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}
func TestPostSecret(t *testing.T) {
	tests := []struct {
		name       string
		host       string
		token      string
		meta       string
		secretType string
		data       interface{}
		id         uint
		status     int
		body       string
		hasErr     bool
	}{
		{
			name:       "successful post secret ByteSlice",
			host:       "", // the host will be replaced with the test server host
			token:      "testToken",
			meta:       "testMeta",
			secretType: "ByteSlice",
			data:       []byte("testData"),
			id:         1,
			status:     http.StatusOK,
			body:       "",
			hasErr:     false,
		},
		{
			name:       "successful post secret not ByteSlice",
			host:       "",
			token:      "testToken",
			meta:       "testMeta",
			secretType: "notByteSlice",
			data:       "testData",
			id:         1,
			status:     http.StatusOK,
			body:       "",
			hasErr:     false,
		},
		{
			name:       "failed post secret due to missing token",
			host:       "",
			token:      "",
			meta:       "testMeta",
			secretType: "ByteSlice",
			data:       []byte("testData"),
			id:         1,
			status:     http.StatusUnauthorized,
			body:       "",
			hasErr:     true,
		},
		{
			name:       "client error",
			host:       "",
			token:      "testToken",
			meta:       "testMeta",
			secretType: "ByteSlice",
			data:       []byte("testData"),
			id:         1,
			status:     http.StatusBadRequest,
			body:       `{"message":"bad request"}`,
			hasErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			err := PostSecret(client, tt.host, tt.token, tt.meta, tt.secretType, tt.data, tt.id)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}

func TestPostTextSecret(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		token  string
		meta   string
		value  string
		id     uint
		status int
		body   string
		hasErr bool
	}{
		{
			name:   "successful post text secret",
			host:   "", // the host will be replaced with the test server host
			token:  "testToken",
			meta:   "testMeta",
			value:  "testValue",
			id:     1,
			status: http.StatusOK,
			body:   "",
			hasErr: false,
		},
		{
			name:   "failed post text secret due to missing token",
			host:   "",
			token:  "",
			meta:   "testMeta",
			value:  "testValue",
			id:     1,
			status: http.StatusUnauthorized,
			body:   "",
			hasErr: true,
		},
		{
			name:   "client error",
			host:   "",
			token:  "testToken",
			meta:   "testMeta",
			value:  "testValue",
			id:     1,
			status: http.StatusBadRequest,
			body:   `{"message":"bad request"}`,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			err := PostTextSecret(client, tt.host, tt.token, tt.meta, tt.value, tt.id)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}

func TestPostKVSecret(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		token  string
		meta   string
		key    string
		value  string
		id     uint
		status int
		body   string
		hasErr bool
	}{
		{
			name:   "successful post KV secret",
			host:   "", // the host will be replaced with the test server host
			token:  "testToken",
			meta:   "testMeta",
			key:    "testKey",
			value:  "testValue",
			id:     1,
			status: http.StatusOK,
			body:   "",
			hasErr: false,
		},
		{
			name:   "failed post KV secret due to missing token",
			host:   "",
			token:  "",
			meta:   "testMeta",
			key:    "testKey",
			value:  "testValue",
			id:     1,
			status: http.StatusUnauthorized,
			body:   "",
			hasErr: true,
		},
		{
			name:   "client error",
			host:   "",
			token:  "testToken",
			meta:   "testMeta",
			key:    "testKey",
			value:  "testValue",
			id:     1,
			status: http.StatusBadRequest,
			body:   `{"message":"bad request"}`,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			err := PostKVSecret(client, tt.host, tt.token, tt.meta, tt.key, tt.value, tt.id)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}

func TestPostCCSecret(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		token   string
		meta    string
		cnn     string
		exp     string
		cvv     string
		cholder string
		id      uint
		status  int
		body    string
		hasErr  bool
	}{
		{
			name:    "successful post CC secret",
			host:    "", // the host will be replaced with the test server host
			token:   "testToken",
			meta:    "testMeta",
			cnn:     "4111111111111111",
			exp:     "12/2024",
			cvv:     "123",
			cholder: "John Doe",
			id:      1,
			status:  http.StatusOK,
			body:    "",
			hasErr:  false,
		},
		{
			name:    "failed post CC secret due to missing token",
			host:    "",
			token:   "",
			meta:    "testMeta",
			cnn:     "4111111111111111",
			exp:     "12/2024",
			cvv:     "123",
			cholder: "John Doe",
			id:      1,
			status:  http.StatusUnauthorized,
			body:    "",
			hasErr:  true,
		},
		{
			name:    "client error",
			host:    "",
			token:   "testToken",
			meta:    "testMeta",
			cnn:     "4111111111111111",
			exp:     "12/2024",
			cvv:     "123",
			cholder: "John Doe",
			id:      1,
			status:  http.StatusBadRequest,
			body:    `{"message":"bad request"}`,
			hasErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			err := PostCCSecret(client, tt.host, tt.token, tt.meta, tt.cnn, tt.exp, tt.cvv, tt.cholder, tt.id)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}

func TestPostFileSecret(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		token  string
		meta   string
		path   string
		id     uint
		status int
		body   string
		hasErr bool
	}{
		{
			name:   "successful post file secret",
			host:   "", // the host will be replaced with the test server host
			token:  "testToken",
			meta:   "testMeta",
			path:   "testFile.txt",
			id:     1,
			status: http.StatusOK,
			body:   "",
			hasErr: false,
		},
		{
			name:   "failed post file secret due to missing token",
			host:   "",
			token:  "",
			meta:   "testMeta",
			path:   "testFile.txt",
			id:     1,
			status: http.StatusUnauthorized,
			body:   "",
			hasErr: true,
		},
		{
			name:   "client error",
			host:   "",
			token:  "testToken",
			meta:   "testMeta",
			path:   "testFile.txt",
			id:     1,
			status: http.StatusBadRequest,
			body:   `{"message":"bad request"}`,
			hasErr: true,
		},
	}

	// Creating a temporary test file
	err := os.WriteFile("testFile.txt", []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("couldn't create a test file: %v", err)
	}
	defer os.Remove("testFile.txt")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			err := PostFileSecret(client, tt.host, tt.token, tt.meta, tt.path, tt.id)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}
func TestDeleteSecret(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		token  string
		id     string
		status int
		body   string
		hasErr bool
	}{
		{
			name:   "successful delete secret",
			host:   "", // the host will be replaced with the test server host
			token:  "testToken",
			id:     "1",
			status: http.StatusOK,
			body:   "",
			hasErr: false,
		},
		{
			name:   "failed delete secret due to missing token",
			host:   "",
			token:  "",
			id:     "1",
			status: http.StatusUnauthorized,
			body:   "",
			hasErr: true,
		},
		{
			name:   "client error",
			host:   "",
			token:  "testToken",
			id:     "1",
			status: http.StatusBadRequest,
			body:   `{"message":"bad request"}`,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintln(w, tt.body)
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			err := DeleteSecret(client, tt.host, tt.token, tt.id)

			if tt.hasErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.hasErr && err != nil {
				t.Fatalf("didn't expect error, got %v", err)
			}
		})
	}
}
func TestGetSecret(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		token       string
		id          string
		status      int
		body        *secret.Secret
		expectedErr string
	}{
		{
			name:   "successful get secret",
			host:   "", // the host will be replaced with the test server host
			token:  "testToken",
			id:     "1",
			status: http.StatusOK,
			body: &secret.Secret{
				ID:         1,
				UserID:     1,
				SecretType: "Text",
				Metadata:   "metadata",
				Value:      []byte("value"),
			},
			expectedErr: "",
		},
		{
			name:        "failed get secret due to missing token",
			host:        "",
			token:       "",
			id:          "1",
			status:      http.StatusUnauthorized,
			body:        nil,
			expectedErr: "unexpected status: 401 Unauthorized",
		},
		{
			name:        "client error",
			host:        "",
			token:       "testToken",
			id:          "1",
			status:      http.StatusBadRequest,
			body:        nil,
			expectedErr: "unexpected status: 400 Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				if tt.body != nil {
					secretJSON, _ := json.Marshal(tt.body)
					fmt.Fprintln(w, string(secretJSON))
				}
			}))
			defer ts.Close()

			// Remove the "http://" from the test server URL and assign it to the host variable
			tt.host = strings.TrimPrefix(ts.URL, "http://")

			client := &http.Client{}
			secret, err := GetSecret(client, tt.host, tt.token, tt.id)

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("expected an error, got nil")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
				return
			}

			if secret.ID != tt.body.ID || string(secret.Value) != string(tt.body.Value) || secret.Metadata != tt.body.Metadata || secret.SecretType != tt.body.SecretType || secret.UserID != tt.body.UserID {
				t.Errorf("expected secret %v, got %v", tt.body, secret)
			}
		})
	}
}
