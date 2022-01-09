package register

import (
	"os"
	"fmt"
	"bytes"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/common"

	"github.com/marcinbor85/nes/api"
)

type RegisterContext struct {
	PublicKeyFile *string
	Email		  *string
}

var Cmd = &cmd.Command{
	Name: "register",
	Help: "Register username at PubKey Service",
	Context: &RegisterContext{},
	Init: Init,
	Logic: Logic,
}

func Init(c *cmd.Command) {
	ctx := c.Context.(*RegisterContext)
	
	ctx.PublicKeyFile = c.Cmd.String("P", "public", &argparse.Options{
		Required: true,
		Help:     `Public key file.`,
	})
	ctx.Email = c.Cmd.String("e", "email", &argparse.Options{
		Required: true,
		Help:     `User email (need to activate username).`,
	})
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*RegisterContext)

	f, e := os.Open(*ctx.PublicKeyFile)
	if e != nil {
		fmt.Println("cannot open public key file")
		return
	}
	defer f.Close()

	bytesBuf := &bytes.Buffer{}
	bytesBuf.ReadFrom(f)
	publicKeyPem := bytesBuf.String()

	pubkeyClient := &api.Client{
		Address: common.G.PubKeyAddress,
	}

	err := pubkeyClient.RegisterNewUsername(common.G.Username, *ctx.Email, publicKeyPem)
	if err != nil {
		fmt.Printf(err.E.Error())
		return
	}
	fmt.Println("username registered. check email for activation.")
}
