package core

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

type TwitchSecret struct {
	Auth     string
	Username string
	Id       string
}

type DiscordSecret struct {
	Auth     string
	Username string
	Id       string
}

type Secrets struct {
	Twitch  TwitchSecret  `toml:"twitch"`
	Discord DiscordSecret `toml:"discord"`
}

type SecretsLoader struct {
	secrets Secrets
}

func NewSecretsLoader(path string) (*SecretsLoader, error) {
	loader := &SecretsLoader{}
	if err := loader.LoadSecrets(path); err != nil {
		return nil, err
	}

	return loader, nil
}

func (s *SecretsLoader) LoadSecrets(path string) error {
	secretsFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(secretsFile), &s.secrets)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretsLoader) GetSecrets() Secrets {
	return s.secrets
}

type TodayData struct {
	Text string
}

type ProjectData struct {
	Text string
}

type Data struct {
	Today   TodayData   `toml:"today"`
	Project ProjectData `toml:"project"`
}

type DataLoader struct {
	data Data
}

func NewDataLoader(path string) (*DataLoader, error) {
	loader := &DataLoader{}
	if err := loader.LoadData(path); err != nil {
		return nil, err
	}

	return loader, nil
}

func (d *DataLoader) LoadData(path string) error {
	dataFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(dataFile), &d.data)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataLoader) WriteData(path string, data *Data) error {
	var buf = new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(&data)
	if err != nil {
		return err
	}

	os.WriteFile(path, buf.Bytes(), 0666)

	return nil
}

func (d *DataLoader) GetData() Data {
	return d.data
}
