package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/structs"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type generateCommand struct {
	cmd.Command
	raw     *bool
	name    *string
	keyPath *string
}

func (c *generateCommand) Register(cmd *kingpin.CmdClause) {

	u, err := user.Current()
	if err != nil {
		panic("Cannot find the current system user.")
	}
	keyPath := path.Join(u.HomeDir, ".ssh")

	command := cmd.Command("generate", "Generate a keypair").Action(c.action)
	c.name = command.Arg("name", "The keypair name").Required().String()
	c.keyPath = command.Flag("key-dir", "Directory to save the keypair to").Default(keyPath).ExistingFileOrDir()
}

func (c *generateCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2046)
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)
	pubKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	publicKey := string(ssh.MarshalAuthorizedKey(pubKey))

	keypair, err := c.Application.APIClient.Keypair().Create(*c.name, publicKey)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		err := ioutil.WriteFile(path.Join(*c.keyPath, "id_"+*c.name), privateKeyBytes, 0600)
		if err != nil {
			return err
		}
		if *c.raw {
			keyPairMap := structs.Map(keypair)
			keypairBytes, _ := json.MarshalIndent(keyPairMap, "", "  ")
			fmt.Println(string(keypairBytes))
		} else {
			logrus.Infof("Keypair '%s' created with an ID of '%s' and saved to '%s'", keypair.Name, keypair.ID, *c.keyPath)
		}
	}

	return nil
}
