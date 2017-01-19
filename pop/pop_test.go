package pop

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/mcilloni/openbaton-docker/pop/server"
	log "github.com/sirupsen/logrus"
)

const (
	uname = "user_name"
	pass  = "pass_value"
	laddr = "localhost:60000"
)

var (
	cfg  server.Config
	conn *grpc.ClientConn
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
		Users:   Users{user.Name: user},
	}

	srv := &server.Server{Config: cfg}

	go func() {
		if err := srv.Serve(); err != nil {
			log.WithError(err).Fatal("Serve failed")
		}
	}()

	conn, err = grpc.Dial(laddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
}

func login() (PopClient, string, error) {
	return loginUserPass(uname, pass)
}

func loginUserPass(uname, pass string) (PopClient, string, error) {
	cln := NewPopClient(conn)
	tok, err := cln.Login(context.Background(), &Credentials{
		Username: uname,
		Password: pass,
	})

	if err != nil {
		return nil, "", err
	}

	if tok == nil {
		return nil, "", errors.New("Invalid response")
	}

	return cln, tok.Value, nil
}

func TestLogin(tst *testing.T) {
	_, _, err := login()
	if err != nil {
		tst.Error(err)
	}
}

func TestLoginFail(tst *testing.T) {
	_, _, err := loginUserPass("wrong user", "random pass")
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
