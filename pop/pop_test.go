package pop

import (
	"context"
	"testing"
	
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/server"
	log "github.com/sirupsen/logrus"
)

const (
	uname = "user_name"
	pass  = "pass_value"
	laddr = "localhost:60000"
)

var (
	cfg server.Config
	cln = client.Client{
		Credentials: client.Credentials{
			Username: uname,
			Password: pass,
		},
	}
)

func init() {
	var err error

	user, e := server.NewUser(uname, pass)
	if e != nil {
		panic(e)
	}

	cfg = server.Config{
		Proto:   "tcp",
		Netaddr: laddr,
		Users:   server.Users{
			user.Name: user,
		},
	}

	srv := &server.Server{Config: cfg}

	go func() {
		if err := srv.Serve(); err != nil {
			log.WithError(err).Fatal("Serve failed")
		}
	}()
}

func TestLogin(tst *testing.T) {
	infos, err := cln.Info(context.Background())
	if err != nil {
		tst.Error(err)
	}

	tst.Log(infos)
}

func TestLoginFail(tst *testing.T) {
	brokenClient := client.Client{
		Credentials: client.Credentials{
			Username: "wrong user", 
			Password: "random pass",
		},
	}
	
	_, err := brokenClient.Info(context.Background())

	if err == nil {
		tst.Error("should have failed")
	}

	tst.Log(err)
}

func TestLogout(tst *testing.T) {
	cln, tok, err := login()
}

func TestUnauthorized(tst *testing.T) {
	cln := NewPopClient(conn)
	_, err := cln.Containers(context.Background(), &Filter{})
	if err == nil {
		tst.Error("should have failed")
	}

	tst.Log(err)
}
