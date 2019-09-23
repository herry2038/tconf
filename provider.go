package tconf

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigFileProvider local file processor
type ConfigFileProvider interface {
	WriteConfig(values map[string]interface{}) error
	ReadConfig() ([]byte, error)
	Unmarshal([]byte) (map[string]interface{}, error)
}

// YAMLProvider yaml processor
type YAMLProvider struct {
	FileName string
}

// WriteConfig write to config file
func (provider YAMLProvider) WriteConfig(values map[string]interface{}) error {
	f, err := os.OpenFile(provider.FileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := yaml.Marshal(values)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return f.Sync()
}

// ReadConfig read data from yaml
func (provider YAMLProvider) ReadConfig() (b []byte, err error) {
	f, err := os.Open(provider.FileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	b = buf.Bytes()

	return
}

// Unmarshal convert byte type to slice type
func (provider YAMLProvider) Unmarshal(bytes []byte) (m map[string]interface{}, err error) {
	err = yaml.Unmarshal(bytes, &m)
	return
}
