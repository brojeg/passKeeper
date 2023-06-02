package cmd

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"path/filepath"
	"strconv"
	"time"

	"passKeeper/internal/cmd/tui/list"
	secret "passKeeper/internal/models/secret"
	client "passKeeper/pkg"
	clientRequest "passKeeper/pkg"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v2"
)

const AppName = "passKeeper"
const cfgFile = "config.yaml"

type config struct {
	Server struct {
		Username string
		Password string
		Token    string
		Host     string
	}
}

type Application struct {
	Config config
	client *http.Client
}

type Username struct {
	Username string `yaml:"username,omitempty"`
}

func GetApplication() *Application {
	var app Application
	token, err := GetKey("token")
	if err != nil {
		log.Fatal(err)
	}
	host, err := GetKey("host")
	if err != nil {
		log.Fatal(err)
	}

	app.Config.Server.Token = token
	app.Config.Server.Host = host

	app.client = &http.Client{Timeout: time.Second * 10}

	return &app
}

func (app Application) Setup() *Application {
	cfg, err := GetUsername()
	if err != nil {
		log.Fatal("Can't read configuration. Try running `passKeeper config` to fix the issue.", "err", err)
	}

	if cfg.Username == "" {
		log.Fatal("Username not set. Run `passKeeper config` first.")
	}
	app.Config.Server.Username = cfg.Username

	app.Config.Server.Password, err = GetKey(AppName)
	if err != nil {
		os.Exit(1)
	}

	app.client = &http.Client{Timeout: time.Second * 10}

	app = *app.Register()
	if app.Config.Server.Token == "" {
		log.Fatal("Registration on server has failed")
	}

	return &app
}
func (app Application) Login() *Application {
	cfg, err := GetUsername()
	if err != nil {
		log.Fatal("Can't read configuration. Try running `passKeeper config` to fix the issue.", "err", err)
	}

	if cfg.Username == "" {
		log.Fatal("Username not set. Run `passKeeper config` first.")
	}
	app.Config.Server.Username = cfg.Username

	app.Config.Server.Password, err = GetKey(AppName)
	if err != nil {
		os.Exit(1)
	}

	app.client = &http.Client{Timeout: time.Second * 10}

	app = *app.doLogin()
	if app.Config.Server.Token == "" {
		log.Fatal("Registration on server has failed")
	}

	return &app
}
func (app Application) doLogin() *Application {
	account, err := clientRequest.SendLoginRequest(app.client, app.Config.Server.Host, app.Config.Server.Username, app.Config.Server.Password)
	if err != nil {
		log.Error(err)
	}
	app.Config.Server.Token = account

	return &app
}
func (app Application) Register() *Application {
	account, err := clientRequest.SendRegisterRequest(app.client, app.Config.Server.Host, app.Config.Server.Username, app.Config.Server.Password)
	if err != nil {
		log.Error(err)
	}
	app.Config.Server.Token = account.Token

	return &app
}

func (app Application) login() *Application {
	cfg, err := GetUsername()
	if err != nil {
		log.Fatal("Can't read configuration. Try running `passKeeper config` to fix the issue.", "err", err)
	}

	if cfg.Username == "" {
		log.Fatal("Username not set. Run `passKeeper config` first.")
	}
	username := cfg.Username

	password, err := GetKey(AppName)
	if err != nil {
		os.Exit(1)
	}
	host, err := GetKey("host")
	if err != nil {
		os.Exit(1)
	}

	token, err := clientRequest.SendLoginRequest(app.client, host, username, password)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	app.Config.Server.Token = token

	return &app
}
func (app Application) ListSecrets() *[]secret.Secret {

	app = *app.login()
	secrets, err := clientRequest.SendGetSecretList(app.client, app.Config.Server.Host, app.Config.Server.Token)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}

	return &secrets
}

func List(app Application, ds []secret.DecodedSecret) {
	columns := []table.Column{
		{Title: "SecretID", Width: 10},
		{Title: "Metadata", Width: 25},
		{Title: "Data", Width: 80},
	}

	var rows []table.Row
	for _, v := range ds {
		meta, err := GetFileInfo(v.Metadata)
		if err != nil {
			rows = append(rows, []string{strconv.Itoa(int(v.ID)), v.Metadata, v.ValueToString()})
		} else {
			filename := []string{meta[0], meta[1]}
			rows = append(rows, []string{strconv.Itoa(int(v.ID)), strings.Join(filename, "."), v.ValueToString()})
		}

	}

	m := list.NewModel(columns, rows)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start passKeeper: %s\n", err)
		os.Exit(1)
	}
}

func (app Application) CreateTextSecret(meta, data string) error {

	app = *app.login()
	err := client.PostTextSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, data, 0)
	if err != nil {
		return err
	}
	return nil

}
func (app Application) EditTextSecret(id uint, meta, data string) error {

	app = *app.login()
	err := client.PostTextSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, data, id)
	if err != nil {
		return err
	}
	return nil

}
func (app Application) CreateKVSecret(meta, key, value string) error {

	app = *app.login()
	err := client.PostKVSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, key, value, 0)
	if err != nil {
		return err
	}
	return nil

}
func (app Application) EditKVSecret(id uint, meta, key, value string) error {

	app = *app.login()
	err := client.PostKVSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, key, value, id)
	if err != nil {
		return err
	}
	return nil

}

func (app Application) CreateCCSecret(meta, cnn, exp, cvv, cholder string) error {

	app = *app.login()
	err := client.PostCCSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, cnn, exp, cvv, cholder, 0)
	if err != nil {
		return err
	}
	return nil

}
func (app Application) CreateFileSecret(meta, path string) error {

	app = *app.login()
	err := client.PostFileSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, path, 0)
	if err != nil {
		return err
	}
	return nil

}

func (app Application) EditCCSecret(id uint, meta, cnn, exp, cvv, cholder string) error {

	app = *app.login()
	err := client.PostCCSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, meta, cnn, exp, cvv, cholder, id)
	if err != nil {
		return err
	}
	return nil

}

func (app Application) GetSecret(id string) (*secret.Secret, error) {

	app = *app.login()
	sec, err := client.GetSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, id)
	if err != nil {
		return nil, err
	}
	return sec, nil

}

func (app Application) DeleteSecret(id string) error {

	app = *app.login()
	err := client.DeleteSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, id)
	if err != nil {
		return err
	}
	return nil

}

func (app Application) DumpSecret(id string) (string, error) {

	app = *app.login()
	sec, err := client.GetSecret(app.client, app.Config.Server.Host, app.Config.Server.Token, id)
	if err != nil {
		return "", err
	}

	if sec.SecretType != "ByteSlice" {
		return "", fmt.Errorf("only bynary data could be saved on disk")
	}
	var secrets []secret.Secret
	secrets = append(secrets, *sec)
	decoded, err := secret.GetDecodedSecrets(secrets)
	if err != nil {
		return "", fmt.Errorf("cannot decode secret. %s", err.Error())
	}

	data := decoded[0].Value.(*secret.ByteSlice)
	path, err := SaveBinarySecretOnDisk(*data, sec.Metadata)
	if err != nil {
		return "", fmt.Errorf("cannot save data on disk. %s", err.Error())
	}

	return path, nil

}

// SetKey creates a new entry in the OS keyring
func SetKey(service, secret string) error {
	err := keyring.Set(service, AppName, secret)
	if err != nil {
		log.Error("failed to fetch secret", "err", err)
		return err
	}

	return nil
}

func DeleteKey(service string) error {
	err := keyring.Delete(service, AppName)
	if err != nil {
		log.Error("failed to delete secret", "err", err)
		return err
	}

	return nil
}

// GetKey retrieves a key from the OS keyring
func GetKey(service string) (string, error) {
	// get password
	secret, err := keyring.Get(service, AppName)
	if err != nil {
		log.Error("failed to fetch key", "err", err)
		return "", err
	}

	return secret, nil
}

// GetUsername from configuration stored on the disk
func GetUsername() (*Username, error) {
	var username Username

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(home, "passKeeper", ".config", cfgFile)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Username{}, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(data, &username)
	if err != nil {
		return nil, err
	}

	return &username, nil
}

func ClearLocalData() error {
	err := DeleteKey(AppName)
	if err != nil {
		return err
	}
	err = DeleteKey("token")
	if err != nil {
		return err
	}
	err = DeleteKey("host")
	if err != nil {
		return err
	}
	err = DeleteFolderStructure()
	if err != nil {
		return err
	}
	return nil
}

func DeleteFolderStructure() error {
	home, _ := os.UserHomeDir()
	fullPath := filepath.Join(home, "passKeeper")

	// clean up after test

	err := os.RemoveAll(fullPath)
	if err != nil {
		return err
	}
	return nil
}

func SaveBinarySecretOnDisk(data []byte, meta string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, "passKeeper", "data")
	log.Print(fullPath)
	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		log.Print("cannot create a file")
		return "", err
	}
	matches, err := GetFileInfo(meta)
	if err != nil {
		return "", err
	}

	fullFilename := matches[0] + "." + matches[1]
	f, err := os.Create(filepath.Join(fullPath, fullFilename))
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return "", err
	}
	return filepath.Join(fullPath, fullFilename), nil

}

func GetFileInfo(meta string) ([]string, error) {
	re := regexp.MustCompile(`^([^|]+)\|([^|]+)\|(.+)$`)
	matches := re.FindStringSubmatch(meta)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid meta string: %s", meta)
	}
	return matches[1:], nil
}
func SetUsername(username string) error {
	creds, err := GetUsername()
	if err != nil {
		return err
	}
	if username != "" {
		creds.Username = username
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(home, "passKeeper", ".config")
	log.Print(fullPath)
	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		log.Print("cannot create a file")
		return err
	}

	f, err := os.OpenFile(filepath.Join(fullPath, cfgFile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&creds)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func PingServer(address string) bool {
	timeout := time.Second * 5
	s := strings.Split(address, ":")
	if len(s) != 2 {
		// address is not in expected format
		return false
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(s[0], s[1]), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}
