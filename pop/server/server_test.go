package server

import (
	"errors"
	"testing"

	"github.com/mcilloni/test-svc/svc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	conn *grpc.ClientConn
)

func init() {
	newConn, err := grpc.Dial("localhost:60000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	conn = newConn
}

func login(uname, pass string) (svc.SvcClient, string, error) {
	cln := svc.NewSvcClient(conn)
	ls, err := cln.Login(context.Background(), &svc.LoginData{
		Username: uname,
		Password: pass,
	})

	if err != nil {
		return nil, "", err
	}

	if ls == nil {
		return nil, "", errors.New("Invalid response")
	}

	return cln, ls.Token, nil
}

func doThings() (string, error) {
	cln, tok, err := login("awawa", "tasty walnuts")
	if err != nil {
		return "", err
	}

	md := metadata.Pairs("token", tok)
	ctx := metadata.NewContext(context.Background(), md)

	repl, err := cln.DoStuff(ctx, &svc.Request{Txt: "client says: banana"})
	if err != nil {
		return "", err
	}

	if repl == nil {
		return "", errors.New("empty response")
	}

	return repl.Txt, nil
}

func TestLogin(tst *testing.T) {
	_, _, err := login("awawa", "tasty walnuts")
	if err != nil {
		tst.Error(err)
	}
}

func TestDoThings(tst *testing.T) {
	txt, err := doThings()
	if err != nil {
		tst.Error(err)
	}

	tst.Log(txt)
}

func TestLoginFail(tst *testing.T) {
	_, _, err := login("awawa", "random pass")
	if err == nil {
		tst.Error("should have failed")
	}

	tst.Log(err)
}

func TestUnauthorized(tst *testing.T) {
	cln := svc.NewSvcClient(conn)
	_, err := cln.DoStuff(context.Background(), &svc.Request{Txt: "client says: banana"})
	if err == nil {
		tst.Error("should have failed")
	}

	tst.Log(err)
}
