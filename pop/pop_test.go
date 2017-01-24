package pop

import (
	"context"
	"encoding/json"
	"testing"

	"google.golang.org/grpc"
	
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
	"github.com/mcilloni/openbaton-docker/pop/server"
	log "github.com/sirupsen/logrus"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

const (
	laddr = "localhost:60000"
	uname = "user_name"
	pass  = "pass_value"
)

var (
	cfg server.Config
	cln = client.Client{
		Credentials: creds.Credentials{
			Host: laddr,
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

func TestFlavours(tst *testing.T) {
	flavs, err := cln.Flavours(context.Background())
	if err != nil {
		tst.Error(err)
	}

	fj, err := json.MarshalIndent(flavs, "", "  ")
	if err != nil {
		tst.Error(err)
	}

	tst.Log(string(fj))
}

func TestImages(tst *testing.T) {
	imgs, err := cln.Images(context.Background())
	if err != nil {
		tst.Error(err)
	}

	ij, err := json.MarshalIndent(imgs, "", "  ")
	if err != nil {
		tst.Error(err)
	}

	tst.Log(string(ij))
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
		Credentials: creds.Credentials{
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

func TestNetworks(tst *testing.T) {
	nets, err := cln.Networks(context.Background())
	if err != nil {
		tst.Error(err)
	}

	nj, err := json.MarshalIndent(nets, "", "  ")
	if err != nil {
		tst.Error(err)
	}

	tst.Log(string(nj))
}

func TestServers(tst *testing.T) {
	srvs, err := cln.Servers(context.Background())
	if err != nil {
		tst.Error(err)
	}

	sj, err := json.MarshalIndent(srvs, "", "  ")
	if err != nil {
		tst.Error(err)
	}

	tst.Log(string(sj))
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
