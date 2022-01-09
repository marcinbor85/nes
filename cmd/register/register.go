package register

import (
	"fmt"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/common"
	"github.com/marcinbor85/nes/crypto/rsa"

	"github.com/marcinbor85/nes/api"
)

type RegisterContext struct {
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
	
	ctx.Email = c.Cmd.String("e", "email", &argparse.Options{
		Required: true,
		Help:     `User email (need to activate username).`,
	})
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*RegisterContext)

	_, publicKeyPem, err := rsa.LoadPublicKey(common.G.PublicKeyFile)
	if err != nil {
		fmt.Println("Cannot load public key:", err.Error())
		return
	}

	pubkeyClient := &api.Client{
		Address: common.G.PubKeyAddress,
	}

	err = pubkeyClient.RegisterNewUsername(common.G.Username, *ctx.Email, publicKeyPem)
	if err != nil {
		fmt.Println("Cannot register username:", err.Error())
		return
	}
	fmt.Println("Username registered. Check email for activation.")
}
