package generate

import (
	"fmt"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/common"

	r "github.com/marcinbor85/nes/crypto/rsa"
)

type GenerateContext struct {
	Size *int
}

var Cmd = &cmd.Command{
	Name: "generate",
	Help: "Generate private and public keys pair",
	Context: &GenerateContext{},
	Init: Init,
	Logic: Logic,
}

func Init(c *cmd.Command) {
	ctx := c.Context.(*GenerateContext)
	ctx.Size = c.Cmd.Int("s", "size", &argparse.Options{Required: false, Help: "Key size", Default: 2048})
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*GenerateContext)
	priv, pub, err := r.GenerateKeysPair(*ctx.Size)
	if err != nil {
		fmt.Println("Generate keys error:", err.Error())
	}
	err = r.SavePrivateKey(priv, common.G.PrivateKeyFile)
	if err != nil {
		fmt.Println("Cannot save private key:", err.Error())
	}
	err = r.SavePublicKey(pub, common.G.PublicKeyFile)
	if err != nil {
		fmt.Println("Cannot save public key:", err.Error())
	}
	fmt.Println("Keys generated successfully.")
}
