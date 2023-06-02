package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	account "passKeeper/internal/models/account"
	secret "passKeeper/internal/models/secret"
)

func sendJSONRequest(client *http.Client, method, host, endpoint, token string, payload interface{}) ([]byte, error) {
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return nil, err
	}

	url := "http://" + host + endpoint
	req, err := http.NewRequest(method, url, payloadBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func SendRegisterRequest(client *http.Client, host, login, password string) (*account.Account, error) {
	if host == "" || login == "" || password == "" {
		return nil, fmt.Errorf("incomplete request")
	}

	data := account.Account{Login: login, Password: password}
	body, err := sendJSONRequest(client, "POST", host, "/api/account/register", "", data)
	if err != nil {
		return nil, err
	}

	var response account.Account
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
func SendGetSecretList(client *http.Client, host, token string) ([]secret.Secret, error) {
	body, err := sendJSONRequest(client, "GET", host, "/api/secrets", token, nil)
	if err != nil {
		return nil, err
	}

	var secrets []secret.Secret
	if err := json.Unmarshal(body, &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

func SendLoginRequest(client *http.Client, host, login, password string) (string, error) {
	if host == "" || login == "" || password == "" {
		return "", fmt.Errorf("incomplete request")
	}

	data := account.Account{Login: login, Password: password}
	body, err := sendJSONRequest(client, "POST", host, "/api/account/login", "", data)
	if err != nil {
		return "", err
	}

	var response string
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response, nil
}

func PostSecret(client *http.Client, host, token, meta, secretType string, data interface{}, id uint) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if secretType == "ByteSlice" {
		_, err = sendJSONRequest(client, "POST", host, "/api/secret", token, data)
		return err
	}

	secretRequest := secret.SecretRequest{ID: id, Type: secretType, Meta: meta, Data: json.RawMessage(dataJson)}
	_, err = sendJSONRequest(client, "POST", host, "/api/secret", token, secretRequest)
	return err
}

func PostTextSecret(client *http.Client, host, token, meta, value string, id uint) error {
	if id == 0 {
		return PostSecret(client, host, token, meta, "Text", secret.Text{Value: value}, 0)
	}
	return PostSecret(client, host, token, meta, "Text", secret.Text{Value: value}, id)
}

func PostKVSecret(client *http.Client, host, token, meta, key, value string, id uint) error {
	if id == 0 {
		return PostSecret(client, host, token, meta, "KeyValue", secret.KeyValue{Key: key, Value: value}, 0)
	}
	return PostSecret(client, host, token, meta, "KeyValue", secret.KeyValue{Key: key, Value: value}, id)
}

func PostCCSecret(client *http.Client, host, token, meta, cnn, exp, cvv, cholder string, id uint) error {
	data := secret.CreditCard{Number: cnn, Expiration: exp, CVV: cvv, Cardholder: cholder}
	if id == 0 {
		return PostSecret(client, host, token, meta, "CreditCard", data, 0)
	}

	return PostSecret(client, host, token, meta, "CreditCard", data, id)
}
func PostFileSecret(client *http.Client, host, token, meta, path string, id uint) error {

	data, err := FiletoBytes(path)
	if err != nil {
		return err
	}

	secret := SecretRequestFromBytes(data, meta)

	return PostSecret(client, host, token, meta, "ByteSlice", secret, id)
}

func DeleteSecret(client *http.Client, host, token, id string) error {
	endpoint := fmt.Sprintf("/api/secret/%s", id)
	_, err := sendJSONRequest(client, "DELETE", host, endpoint, token, nil)
	return err
}

func GetSecret(client *http.Client, host, token, id string) (*secret.Secret, error) {
	endpoint := fmt.Sprintf("/api/secret/%s", id)
	body, err := sendJSONRequest(client, "GET", host, endpoint, token, nil)
	if err != nil {

		return nil, err
	}

	var secretResult secret.Secret
	if err := json.Unmarshal(body, &secretResult); err != nil {

		return nil, err
	}

	return &secretResult, nil
}
