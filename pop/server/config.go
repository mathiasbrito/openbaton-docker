package server

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/BurntSushi/toml" // because it implements TOML Marshalling
	"github.com/docker/docker/client"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

// DefaultConfig is a sane template config for a local server.
var (
	DefaultConfig = Config{
		PopName:     "docker-popd",
		Proto:       pop.DefaultListenProtocol,
		Netaddr:     pop.DefaultListenAddress,
		Users:       Users{},
		DockerdHost: client.DefaultDockerHost,
	}

	ErrMalformedAuthVar = errors.New("invalid POPD_AUTH variable (must be formatted like `user:pass,user2,pass2,[...]`)")
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

// LoadConfig either reads a config file through viper, or
// it initialialises a new config bases on the DefaultConfig template, using the POPD_AUTH config variable
// to fill the authentication data.
// If neither of these options suits you, just fill a Config structure by yourself and create a Server instance using it.
func LoadConfig() (cfg Config, err error) {
	if viper.ConfigFileUsed() != "" {
		err = viper.Unmarshal(&cfg)
		return
	}

	authData := os.Getenv("POPD_AUTH")
	if authData == "" {
		err = errors.New("neither a config file nor POPD_AUTH env variable is set")
		return
	}

	creds := strings.Split(authData, ",")
	err = ErrMalformedAuthVar

	if len(creds) < 0 {
		return
	}

	users := Users{}

	for _, cred := range creds {
		splitCred := strings.Split(cred, ":")
		if len(splitCred) != 2 {
			return
		}

		uname, pass := splitCred[0], splitCred[1]

		user, e := NewUser(uname, pass)
		if e != nil {
			err = e
			return
		}

		users[uname] = user
	}

	cfg = DefaultConfig
	cfg.Users = users

	return cfg, nil
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
