package server

import (
	"io"
	"os"

	"github.com/spf13/viper"
	"github.com/BurntSushi/toml"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/docker/docker/client"
)

// DefaultConfig is a sane template config for a local server.
var DefaultConfig = Config{
	PopName: "docker-popd",
	Proto: pop.DefaultListenProtocol,
	Netaddr: pop.DefaultListenAddress,
	Users: Users{},
	DockerdHost: client.DefaultDockerHost,
}

// Config for the PoP service.
type Config struct {
	PopName     string
	Proto       string
	Netaddr     string
	Users       Users
	DockerdHost string
	// Add TLS here
}

// Load reads a config from viper.
func LoadConfig() (cfg Config, err error) {
	err = viper.Unmarshal(&cfg)
	return
}

func (cfg Config) Store(w io.Writer) error {
	return toml.NewEncoder(w).Encode(cfg)
}

// Store stores a config into a file. It won't replace an existing file if
// overwrite is false.
func (cfg Config) StoreFile(fileName string, overwrite bool) error {
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

	return cfg.Store(file)
}
