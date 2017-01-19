package server

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// Config for the PoP service.
type Config struct {
	PopName     string
	Proto       string
	Netaddr     string
	Users       Users
	DockerdHost string
	// Add TLS here
}

// ReadConfig reads a config from a Reader.
func ReadConfig(r io.Reader) (cfg Config, err error) {
	_, err = toml.DecodeReader(r, &cfg)
	return
}

// ReadConfigFile reads a config from a file.
func ReadConfigFile(fileName string) (cfg Config, err error) {
	_, err = toml.DecodeFile(fileName, &cfg)
	return
}

// Store stores a config into a file. It won't replace an existing file if
// overwrite is false.
func (cfg Config) Store(fileName string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return err
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := toml.NewEncoder(file)

	return enc.Encode(cfg)
}
