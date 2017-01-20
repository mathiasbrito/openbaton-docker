package pop

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/server"
	log "github.com/sirupsen/logrus"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

const (
	uname = "user_name"
	pass  = "pass_value"
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
	user, err := server.NewUser(uname, pass)
	if err != nil {
		panic(err)
	}

	cfg = server.Config{
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
	err := client.FlushSessions()
	if err != nil {
		tst.Error(err)	
	}
}

func TestUnauthorized(tst *testing.T) {
	conn, err := grpc.Dial(laddr, grpc.WithInsecure())
	if err != nil {
		tst.Error(err)
	}
	
	cln := pop.NewPopClient(conn)
	_, err = cln.Containers(context.Background(), &pop.Filter{})
	if err == nil {
		tst.Error("should have failed")
	}

	tst.Log(err)
}
