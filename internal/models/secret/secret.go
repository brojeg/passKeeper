package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
)

type SecretRequest struct {
	ID       uint            `json:"id"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
	ByteData string          `json:"byteData,omitempty"` // New field for base64 encoded []byte
	Meta     string          `json:"meta,omitempty"`
}

type Secret struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint
	Value      ByteSlice
	SecretType string
	Metadata   string
}
type DecodedSecret struct {
	ID       uint
	UserID   uint
	Value    interface{}
	Metadata string
}

func NewSecret(userID uint, secretType string, value ByteConvertible, meta string) (Secret, error) {
	bytes, err := value.ToBytes()
	if err != nil {
		return Secret{}, err
	}

	return Secret{UserID: userID, Value: bytes, SecretType: secretType, Metadata: meta}, nil
}

func GetSecretFromRequest(req SecretRequest, user uint) (ByteConvertible, error) {
	var value ByteConvertible
	switch req.Type {
	case "KeyValue":
		value = &KeyValue{}
	case "Text":
		value = &Text{}
	case "CreditCard":
		value = &CreditCard{}
	case "ByteSlice":
		value = &ByteSlice{}

	default:
		return nil, fmt.Errorf("invalid type: %s", req.Type)
	}

	// Decode ByteData if present
	if req.Type == "ByteSlice" {
		if string(req.ByteData) != "" {
			data := ByteSlice(string(req.ByteData))
			*value.(*ByteSlice) = data
		}
	} else {
		// Otherwise, unmarshal the JSON into the value
		if err := json.Unmarshal(req.Data, value); err != nil {
			return nil, err
		}
	}

	return value, nil
}

func GetDecodedSecrets(secrets []Secret) ([]DecodedSecret, error) {
	decodedSecrets := make([]DecodedSecret, len(secrets))
	for i, secret := range secrets {
		var value ByteConvertible
		switch secret.SecretType {
		case "KeyValue":
			value = new(KeyValue)
		case "Text":
			value = new(Text)
		case "CreditCard":
			value = new(CreditCard)
		case "ByteSlice":
			value = new(ByteSlice)
		default:
			return nil, fmt.Errorf("unknown secret type: %s", secret.SecretType)
		}

		err := value.FromBytes(secret.Value)
		if err != nil {
			return nil, fmt.Errorf("cannot decode secret %d of type %s: %w", secret.ID, secret.SecretType, err)
		}
		decodedSecrets[i] = DecodedSecret{
			ID:       secret.ID,
			UserID:   secret.UserID,
			Value:    value,
			Metadata: secret.Metadata,
		}
	}
	return decodedSecrets, nil
}

func (ds *DecodedSecret) ValueToString() string {
	switch v := ds.Value.(type) {
	case *KeyValue:
		return fmt.Sprintf("Key: %s,\nValue: %s", v.Key, v.Value)
	case *Text:
		return v.Value
	case *CreditCard:
		return fmt.Sprintf("Number: %s,\n Expiration: %s,\n CVV: %s,\n Cardholder: %s", v.Number, v.Expiration, v.CVV, v.Cardholder)
	case *ByteSlice:
		re := regexp.MustCompile(`^([^|]+)\|([^|]+)\|(.+)$`)
		matches := re.FindStringSubmatch(ds.Metadata)
		if len(matches) != 4 {
			return "*Binary data*"
		}
		return matches[3]
	default:
		return "Unknown Value Type"
	}
}

type ByteConvertible interface {
	ToBytes() (ByteSlice, error)
	FromBytes(ByteSlice) error
}

type KeyValue struct {
	Key   string
	Value string
}
type CreditCard struct {
	Number     string
	Expiration string
	CVV        string
	Cardholder string
}
type Text struct {
	Value string
}
type ByteSlice []byte

func (kv *KeyValue) ToBytes() (ByteSlice, error) {
	data, err := json.Marshal(kv)
	if err != nil {
		return nil, err
	}
	return ByteSlice(data), nil
}

func (kv *KeyValue) FromBytes(data ByteSlice) error {
	return json.Unmarshal([]byte(data), kv)
}

func (t *Text) ToBytes() (ByteSlice, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return ByteSlice(data), nil
}

func (t *Text) FromBytes(data ByteSlice) error {
	return json.Unmarshal([]byte(data), t)
}

func (b ByteSlice) ToBytes() (ByteSlice, error) {
	return b, nil
}

func (b *ByteSlice) FromBytes(data ByteSlice) error {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	*b = ByteSlice(decoded)
	return nil
}

func (cc *CreditCard) ToBytes() (ByteSlice, error) {
	data, err := json.Marshal(cc)
	if err != nil {
		return nil, err
	}
	return ByteSlice(data), nil
}

func (cc *CreditCard) FromBytes(data ByteSlice) error {
	return json.Unmarshal([]byte(data), cc)
}

func (t *Text) String() string {
	return fmt.Sprintf("{ \"Value\": \"%s\" }", t.Value)
}

func (kv *KeyValue) String() string {
	return fmt.Sprintf("{ \"Key\": \"%s\", \"Value\": \"%s\" }", kv.Key, kv.Value)
}

func (cc *CreditCard) String() string {
	return fmt.Sprintf("{ \"Number\": \"%s\", \"Expiration\": \"%s\", \"CVV\": \"%s\", \"Cardholder\": \"%s\" }", cc.Number, cc.Expiration, cc.CVV, cc.Cardholder)
}

func (b *ByteSlice) String() string {
	return fmt.Sprintf("{ \"data\": \"%s\" }", base64.StdEncoding.EncodeToString(*b))
}
